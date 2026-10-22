// policy_compiler compiles acp_core policies into Alloy models.
//
// Usage:
//
//	policy_compiler compile <policy.yaml> <output dir>
//	policy_compiler serve [address]
package main

import (
	"context"
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
	fmt.Fprintf(os.Stderr, "usage:\n  %[1]v compile <policy.yaml> <output dir>\n  %[1]v serve [address]\n", filepath.Base(os.Args[0]))
	os.Exit(2)
}

func compile(policyPath, outDir string) {
	definition, err := os.ReadFile(policyPath)
	if err != nil {
		log.Fatal(err)
	}
	// modules are named after the policy file, eg. policy.yaml emits policy.als
	module := strings.TrimSuffix(filepath.Base(policyPath), filepath.Ext(policyPath))
	resp, err := services.NewPolicyCompilerService().Compile(context.Background(), &types.CompileRequest{
		Policy: string(definition),
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
