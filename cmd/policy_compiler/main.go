// policy_compiler compiles acp_core policies into Alloy models.
//
// Usage:
//
//	policy_compiler compile <policy.yaml> <output dir>
//	policy_compiler instance [-o dir] <policy.yaml> <instance.xml>...
//	policy_compiler serve [address]
//
// instance translates instances written by `alloy exec -t xml` back into relationships
// and reports what every actor can do, as evaluated by the acp engine.
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
		if len(os.Args) != 4 {
			usage()
		}
		compile(os.Args[2], os.Args[3])
	case "instance":
		fs := flag.NewFlagSet("instance", flag.ExitOnError)
		outDir := fs.String("o", "", "write each instance as relationships.txt and theorems.txt under this directory")
		fs.Usage = usage
		fs.Parse(os.Args[2:])
		if fs.NArg() < 2 {
			usage()
		}
		instances(fs.Arg(0), *outDir, fs.Args()[1:])
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
	fmt.Fprintf(os.Stderr, "usage:\n  %[1]v compile <policy.yaml> <output dir>\n  %[1]v instance [-o dir] <policy.yaml> <instance.xml>...\n  %[1]v serve [address]\n", filepath.Base(os.Args[0]))
	os.Exit(2)
}

func read(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func compile(policyPath, outDir string) {
	definition := read(policyPath)
	// modules are named after the policy file, eg. policy.yaml emits policy.als
	module := strings.TrimSuffix(filepath.Base(policyPath), filepath.Ext(policyPath))
	resp, err := services.NewPolicyCompilerService().Compile(context.Background(), &types.CompileRequest{
		Policy: definition,
		Module: module,
	})
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
