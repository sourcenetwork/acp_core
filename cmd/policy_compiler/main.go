// policy_compiler compiles acp_core policies into Alloy models.
//
// Usage:
//
//	policy_compiler compile [-a assumptions.yaml] <policy.yaml> <output dir>
//	policy_compiler instance [-o dir] <policy.yaml> <instance.xml>...
//	policy_compiler verify [-a assumptions.yaml] [-s scope] [-o dir] <policy.yaml> <contracts.yaml>
//	policy_compiler serve [address]
//
// instance translates instances written by `alloy exec -t xml` back into relationships
// and reports what every actor can do, as evaluated by the acp engine.
//
// verify checks contracts against every state allowed by the assumptions, up to a scope,
// reporting counterexamples. It exits with 1 unless every contract is satisfied.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"

	"github.com/sourcenetwork/acp_core/internal/policy_compiler"
	"github.com/sourcenetwork/acp_core/pkg/services"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

const defaultAddress = "localhost:8090"

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "compile":
		fs := flag.NewFlagSet("compile", flag.ExitOnError)
		assumptionsPath := fs.String("a", "", "application assumptions yaml, emitted as <module>_assumptions.als")
		fs.Usage = usage
		fs.Parse(os.Args[2:])
		if fs.NArg() != 2 {
			usage()
		}
		compile(fs.Arg(0), fs.Arg(1), *assumptionsPath)
	case "instance":
		fs := flag.NewFlagSet("instance", flag.ExitOnError)
		outDir := fs.String("o", "", "write each instance as relationships.txt and theorems.txt under this directory")
		fs.Usage = usage
		fs.Parse(os.Args[2:])
		if fs.NArg() < 2 {
			usage()
		}
		instances(fs.Arg(0), *outDir, fs.Args()[1:])
	case "verify":
		fs := flag.NewFlagSet("verify", flag.ExitOnError)
		assumptionsPath := fs.String("a", "", "application assumptions yaml")
		scope := fs.Uint("s", 3, "number of objects of each resource, and of actors, in the states checked")
		outDir := fs.String("o", "", "write each counterexample or example as <dir>/<contract>/relationships.txt")
		fs.Usage = usage
		fs.Parse(os.Args[2:])
		if fs.NArg() != 2 {
			usage()
		}
		verify(fs.Arg(0), fs.Arg(1), *assumptionsPath, *outDir, *scope)
	case "serve":
		addr := defaultAddress
		if len(os.Args) > 2 {
			addr = os.Args[2]
		}
		serve(addr)
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage:\n  %[1]v compile [-a assumptions.yaml] <policy.yaml> <output dir>\n  %[1]v instance [-o dir] <policy.yaml> <instance.xml>...\n  %[1]v verify [-a assumptions.yaml] [-s scope] [-o dir] <policy.yaml> <contracts.yaml>\n  %[1]v serve [address]\n", filepath.Base(os.Args[0]))
	os.Exit(2)
}

func read(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func compile(policyPath, outDir, assumptionsPath string) {
	definition := read(policyPath)
	// modules are named after the policy file, eg. policy.yaml emits policy.als
	module := strings.TrimSuffix(filepath.Base(policyPath), filepath.Ext(policyPath))
	req := &types.CompileRequest{
		Policy: definition,
		Module: module,
	}
	if assumptionsPath != "" {
		assumptions, err := policy_compiler.UnmarshalAssumptions(read(assumptionsPath))
		if err != nil {
			log.Fatal(err)
		}
		req.Constraints = assumptions
	}
	resp, err := services.NewPolicyCompilerService().Compile(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatal(err)
	}
	for _, f := range resp.Files {
		path := filepath.Join(outDir, f.Name)
		if err := os.WriteFile(path, []byte(f.Content), 0644); err != nil {
			log.Fatal(err)
		}
		fmt.Println(path)
	}
}

func serve(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	types.RegisterPolicyCompilerServer(server, services.NewPolicyCompilerService())
	log.Printf("policy compiler listening on %v", listener.Addr())
	log.Fatal(server.Serve(listener))
}
