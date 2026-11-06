// acp works with acp_core workspaces: a directory with a policy, its system constraints, contracts and saved worlds
// (see internal/workspace). It runs nothing locally: every operation is a call to the acp service.
//
// Usage:
//
//	acp verify [-s scope] [-o dir] <workspace>
//	acp generate [-n count] [-s scope] [-a actors] [-o dir] <workspace>
//	acp model [-o dir] <workspace>
//	acp check <workspace> <world> <object> <permission> <actor>
//	acp explain <workspace> <world> <object> <permission> <actor>
//	acp playground [-addr address] [-ui dir] [-open=false] <workspace>
//
// playground serves the playground UI for the workspace, reloading it when its files change, and proxies the UI's
// API calls to the acp service.
//
// The service address comes from -server, or the ACP_SERVER environment variable (default http://localhost:8090).
// Exit codes: 0 success (contracts satisfied, access allowed), 1 failure (contracts not satisfied, access denied),
// 2 error.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sourcenetwork/acp_core/internal/playground2"
	"github.com/sourcenetwork/acp_core/internal/workspace"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

const defaultServer = "http://localhost:8090"

// module names the compiled modules, after the workspace's policy.yaml
const module = "policy"

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		usage()
	}
	server := defaultServer
	if env := os.Getenv("ACP_SERVER"); env != "" {
		server = env
	}

	fs := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	fs.Usage = usage
	fs.StringVar(&server, "server", server, "acp service address")
	args := func(n int) []string {
		if fs.NArg() != n {
			usage()
		}
		return fs.Args()
	}
	switch os.Args[1] {
	case "verify":
		scope := fs.Uint("s", 3, "number of objects of each resource, and of actors, in a world")
		outDir := fs.String("o", "", "write each counterexample or example as <dir>/<contract>-<property>.yaml")
		fs.Parse(os.Args[2:])
		a := args(1)
		os.Exit(verify(newClient(server), load(a[0]), *scope, *outDir))
	case "generate":
		count := fs.Uint("n", 1, "maximum number of worlds")
		scope := fs.Uint("s", 3, "number of objects of each resource in a world")
		actors := fs.Uint("a", 3, "number of actors in each world")
		outDir := fs.String("o", "", "write each world as <dir>/generated-<n>.yaml")
		fs.Parse(os.Args[2:])
		a := args(1)
		generate(newClient(server), load(a[0]), *count, *scope, *actors, *outDir)
	case "model":
		outDir := fs.String("o", "", "write the modules to this directory instead of printing them")
		fs.Parse(os.Args[2:])
		a := args(1)
		model(newClient(server), load(a[0]), *outDir)
	case "playground":
		addr := fs.String("addr", "localhost:0", "listen address (default: a free port)")
		uiDir := fs.String("ui", "", "serve the UI from this directory (eg. playground2/dist) instead of the embedded one, eg. while developing it")
		open := fs.Bool("open", true, "open the playground in a browser")
		fs.Parse(os.Args[2:])
		a := args(1)
		playground(server, a[0], *addr, *uiDir, *open)
	case "check", "explain":
		fs.Parse(os.Args[2:])
		a := args(5)
		ws := load(a[0])
		scenario := &types.Scenario{SystemModel: ws.SystemModel, World: findWorld(ws, a[1])}
		object, actor := parseObject(a[2]), types.NewActor(a[4])
		if os.Args[1] == "check" {
			os.Exit(check(newClient(server), scenario, object, a[3], actor))
		}
		os.Exit(explain(newClient(server), scenario, object, a[3], actor))
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `usage:
  acp verify [-s scope] [-o dir] <workspace>
  acp generate [-n count] [-s scope] [-a actors] [-o dir] <workspace>
  acp model [-o dir] <workspace>
  acp check <workspace> <world> <object> <permission> <actor>
  acp explain <workspace> <world> <object> <permission> <actor>
  acp playground [-addr address] [-ui dir] [-open=false] <workspace>

A workspace is a directory with policy.yaml, and optionally constraints.yaml, contracts.yaml and worlds/*.yaml.
<world> names a saved world (worlds/<world>.yaml), <object> is resource:id (doc:readme), <actor> an actor id.
Every command takes -server; the service address defaults to $ACP_SERVER, or %v.
`, defaultServer)
	os.Exit(2)
}

func fatal(v ...any) {
	log.Print(v...)
	os.Exit(2)
}

func load(dir string) *types.Workspace {
	ws, err := workspace.Load(dir)
	if err != nil {
		fatal(err)
	}
	return ws
}

func findWorld(ws *types.Workspace, name string) *types.World {
	var names []string
	for _, w := range ws.Worlds {
		if w.Name == name {
			return w
		}
		names = append(names, w.Name)
	}
	if len(names) == 0 {
		fatal(fmt.Sprintf("no world %q: the workspace has no saved worlds (%v/*.yaml)", name, workspace.WorldsDir))
	}
	fatal(fmt.Sprintf("no world %q, the workspace has: %v", name, strings.Join(names, ", ")))
	return nil
}

func parseObject(s string) *types.Object {
	resource, id, ok := strings.Cut(s, ":")
	if !ok || resource == "" || id == "" {
		fatal(fmt.Sprintf("invalid object %q: expected resource:id, eg. doc:readme", s))
	}
	return types.NewObject(resource, id)
}

// verify prints a report of the workspace's contracts, returning the exit code: 0 if every contract is satisfied,
// drafts aside.
func verify(c *client, ws *types.Workspace, scope uint, outDir string) int {
	if len(ws.SystemModel.Requirements) == 0 {
		fatal(fmt.Sprintf("no requirements: add them to %v", workspace.ContractsFile))
	}
	resp, err := c.Verify(context.Background(), &types.VerifyRequest{
		SystemModel: ws.SystemModel,
		Config:      &types.VerifyConfig{Scope: uint32(scope)},
	})
	if err != nil {
		fatal(err)
	}
	var save func(string, *types.World)
	if outDir != "" {
		save = func(name string, w *types.World) { writeWorld(outDir, name, w) }
	}
	if playground2.PrintReport(os.Stdout, ws.SystemModel.Requirements, resp, save) > 0 {
		return 1
	}
	return 0
}

func generate(c *client, ws *types.Workspace, count, scope, actors uint, outDir string) {
	resp, err := c.Generate(context.Background(), &types.GenerateRequest{
		SystemModel: ws.SystemModel,
		Config:      &types.GenerateConfig{Count: uint32(count), Scope: uint32(scope), Actors: uint32(actors)},
	})
	if err != nil {
		fatal(err)
	}
	if len(resp.Worlds) == 0 {
		fmt.Println("no worlds: the policy and constraints allow none at this scope")
		return
	}
	for i, w := range resp.Worlds {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("world %v\n", i)
		printRelationships(w)
		if outDir != "" {
			writeWorld(outDir, fmt.Sprintf("generated-%v", i), w)
		}
	}
}

// model prints the Alloy modules the workspace's policy and constraints compile to, or writes them to outDir
func model(c *client, ws *types.Workspace, outDir string) {
	resp, err := c.Compile(context.Background(), &types.CompileRequest{
		Policy:      ws.SystemModel.Policy,
		MarshalType: types.PolicyMarshalingType_YAML,
		Module:      module,
		Constraints: ws.SystemModel.Constraints,
	})
	if err != nil {
		fatal(err)
	}
	for i, f := range resp.Files {
		if outDir == "" {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("==> %v <==\n%v", f.Name, f.Content)
			continue
		}
		if err := os.MkdirAll(outDir, 0755); err != nil {
			fatal(err)
		}
		path := filepath.Join(outDir, f.Name)
		if err := os.WriteFile(path, []byte(f.Content), 0644); err != nil {
			fatal(err)
		}
		fmt.Println(path)
	}
}

// check prints whether actor has permission on object, returning 0 if allowed and 1 if not
func check(c *client, scenario *types.Scenario, object *types.Object, permission string, actor *types.Actor) int {
	resp, err := c.Check(context.Background(), &types.CheckRequest{
		Scenario: scenario, Object: object, Permission: permission, Actor: actor,
	})
	if err != nil {
		fatal(err)
	}
	return printDecision(resp.Allowed, object, permission, actor)
}

// explain prints whether actor has permission on object and the engine's search, returning 0 if allowed and 1 if not
func explain(c *client, scenario *types.Scenario, object *types.Object, permission string, actor *types.Actor) int {
	resp, err := c.Explain(context.Background(), &types.ExplainRequest{
		Scenario: scenario, Object: object, Permission: permission, Actor: actor,
	})
	if err != nil {
		fatal(err)
	}
	code := printDecision(resp.Allowed, object, permission, actor)
	fmt.Println()
	printGraph(os.Stdout, resp.Graph)
	return code
}

func printDecision(allowed bool, object *types.Object, permission string, actor *types.Actor) int {
	decision, code := "allowed", 0
	if !allowed {
		decision, code = "denied", 1
	}
	fmt.Printf("%v: %v %v %v:%v\n", decision, actor.Id, permission, object.Resource, object.Id)
	return code
}

func printRelationships(w *types.World) {
	for _, line := range strings.Split(strings.TrimSpace(w.Relationships), "\n") {
		fmt.Printf("    %v\n", line)
	}
}

// writeWorld saves a world as <dir>/<name>.yaml, never overwriting an existing world
func writeWorld(dir, name string, w *types.World) {
	data, err := workspace.MarshalWorld(w)
	if err != nil {
		fatal(err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		fatal(err)
	}
	path := filepath.Join(dir, name+".yaml")
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if os.IsExist(err) {
		fatal(fmt.Sprintf("%v already exists: remove it or choose another -o directory", path))
	}
	if err != nil {
		fatal(err)
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		fatal(err)
	}
	fmt.Printf("         saved as %v\n", path)
}
