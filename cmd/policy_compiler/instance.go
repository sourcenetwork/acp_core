package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sourcenetwork/acp_core/internal/alloy"
	pc "github.com/sourcenetwork/acp_core/internal/policy_compiler"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/services"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// instances translates Alloy XML instances into relationships, evaluates them with the acp engine
// and reports what every actor can do. If outDir is set, each instance is also written as
// relationships.txt and theorems.txt, loadable by acp_simulate and the playground.
func instances(policyPath, outDir string, instancePaths []string) {
	definition := read(policyPath)
	model, err := pc.Compile(definition, types.PolicyMarshalingType_YAML)
	if err != nil {
		log.Fatal(err)
	}

	for i, path := range instancePaths {
		inst, err := alloy.ParseInstance(model, []byte(read(path)))
		if err != nil {
			log.Fatalf("%v: %v", path, err)
		}
		perms, manages := evaluate(definition, model, inst)

		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("== %v (%v)\n", filepath.Base(path), inst.Command)
		printInstance(inst, perms, manages)

		if outDir != "" {
			dir := filepath.Join(outDir, strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
			writeScenario(dir, inst, perms, manages)
			fmt.Printf("\nwritten to %v\n", dir)
		}
	}
}

// grants maps actor -> object -> names of the permissions (or managed relations) the actor has
type grants map[string]map[string][]string

func (g grants) add(actor, object, name string) {
	if g[actor] == nil {
		g[actor] = make(map[string][]string)
	}
	g[actor][object] = append(g[actor][object], name)
}

// evaluate asks the acp engine for every permission and management right of every actor on every object
func evaluate(definition string, model *pc.Model, inst *alloy.Instance) (perms, manages grants) {
	authorizations := &strings.Builder{}
	delegations := &strings.Builder{}
	for _, object := range inst.Objects {
		r := resourceOf(model, object)
		for _, actor := range inst.Actors {
			for _, p := range r.Permissions {
				fmt.Fprintf(authorizations, "  %v#%v@%v\n", object, p.Name, actor)
			}
			for _, a := range r.Authorities {
				fmt.Fprintf(delegations, "  %v > %v#%v\n", actor, object, a.Name)
			}
		}
	}

	manager, err := runtime.NewRuntimeManager()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := services.NewSimulatorService(manager).Evaluate(context.Background(), &types.EvaluateRequest{
		Policy:        definition,
		Relationships: strings.Join(inst.Relationships, "\n"),
		Theorems:      fmt.Sprintf("Authorizations {\n%v}\nDelegations {\n%v}\n", authorizations, delegations),
	})
	if err != nil {
		log.Fatal(err)
	}
	if !resp.Valid {
		log.Fatalf("engine rejected the instance: %v", resp.Errors)
	}

	perms, manages = make(grants), make(grants)
	for _, res := range resp.Result.AuthorizationTheoremsResult {
		thm := res.Result.Theorem
		if res.Result.Result.Status == types.ResultStatus_Accept {
			perms.add(thm.Actor.Id, objectID(thm.Operation.Object), thm.Operation.Permission)
		}
	}
	for _, res := range resp.Result.DelegationTheoremsResult {
		thm := res.Result.Theorem
		if res.Result.Result.Status == types.ResultStatus_Accept {
			manages.add(thm.Actor.Id, objectID(thm.Operation.Object), thm.Operation.Permission)
		}
	}
	return perms, manages
}

func printInstance(inst *alloy.Instance, perms, manages grants) {
	var skolems []string
	for name, values := range inst.Skolems {
		skolems = append(skolems, fmt.Sprintf("%v = %v", name, strings.Join(values, ", ")))
	}
	sort.Strings(skolems)
	for _, s := range skolems {
		fmt.Printf("skolem %v\n", s)
	}

	fmt.Println("\nrelationships:")
	for _, rel := range inst.Relationships {
		fmt.Printf("  %v\n", rel)
	}

	fmt.Println("\nwhat each actor can do:")
	width := 0
	for _, object := range inst.Objects {
		width = max(width, len(object))
	}
	for _, actor := range inst.Actors {
		fmt.Printf("  %v\n", actor)
		any := false
		for _, object := range inst.Objects {
			p, m := perms[actor][object], manages[actor][object]
			if len(p) == 0 && len(m) == 0 {
				continue
			}
			any = true
			line := fmt.Sprintf("    %-*v  %v", width, object, strings.Join(p, " "))
			if len(m) > 0 {
				line += "  | manages " + strings.Join(m, " ")
			}
			fmt.Println(strings.TrimRight(line, " "))
		}
		if !any {
			fmt.Println("    nothing")
		}
	}
}

// writeScenario writes the instance relationships, and theorems asserting what every actor can do
func writeScenario(dir string, inst *alloy.Instance, perms, manages grants) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal(err)
	}
	relationships := strings.Join(inst.Relationships, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(dir, "relationships.txt"), []byte(relationships), 0644); err != nil {
		log.Fatal(err)
	}

	b := &strings.Builder{}
	b.WriteString("Authorizations {\n")
	for _, actor := range inst.Actors {
		for _, object := range inst.Objects {
			for _, p := range perms[actor][object] {
				fmt.Fprintf(b, "  %v#%v@%v\n", object, p, actor)
			}
		}
	}
	b.WriteString("}\n\nDelegations {\n")
	for _, actor := range inst.Actors {
		for _, object := range inst.Objects {
			for _, m := range manages[actor][object] {
				fmt.Fprintf(b, "  %v > %v#%v\n", actor, object, m)
			}
		}
	}
	b.WriteString("}\n")
	if err := os.WriteFile(filepath.Join(dir, "theorems.txt"), []byte(b.String()), 0644); err != nil {
		log.Fatal(err)
	}
}

func resourceOf(model *pc.Model, object string) *pc.Resource {
	name, _, _ := strings.Cut(object, ":")
	for _, r := range model.Resources {
		if r.Name == name {
			return r
		}
	}
	log.Fatalf("unknown resource for object %v", object)
	return nil
}

func objectID(obj *types.Object) string {
	return obj.Resource + ":" + obj.Id
}
