// acp_simulate evaluates theorems against a policy and a set of relationships.
//
// Usage:
//
//	acp_simulate [-json] <policy.yaml> <relationships.txt> [theorems.txt]
//
// Exits with 0 if every theorem holds, 1 if some theorem failed
// and 2 if the inputs couldn't be loaded.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosmos/gogoproto/jsonpb"

	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/services"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func main() {
	asJSON := flag.Bool("json", false, "print the full evaluation response as JSON")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %v [-json] <policy.yaml> <relationships.txt> [theorems.txt]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 2 || flag.NArg() > 3 {
		flag.Usage()
		os.Exit(2)
	}

	req := &types.EvaluateRequest{
		Policy:        read(flag.Arg(0)),
		Relationships: read(flag.Arg(1)),
	}
	if flag.NArg() == 3 {
		req.Theorems = read(flag.Arg(2))
	}

	manager, err := runtime.NewRuntimeManager()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := services.NewSimulatorService(manager).Evaluate(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	if *asJSON {
		marshaler := jsonpb.Marshaler{Indent: "  "}
		if err := marshaler.Marshal(os.Stdout, resp); err != nil {
			log.Fatal(err)
		}
		fmt.Println()
	} else {
		report(req, resp)
	}

	switch {
	case !resp.Valid:
		os.Exit(2)
	case !resp.Result.Ok:
		os.Exit(1)
	}
}

func read(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

// report prints errors or theorem results, quoting each theorem as written in the theorems file
func report(req *types.EvaluateRequest, resp *types.EvaluateResponse) {
	if !resp.Valid {
		fmt.Println("INVALID")
		inputs := []struct {
			name string
			errs []*types.LocatedMessage
		}{
			{"policy", resp.Errors.PolicyErrors},
			{"relationships", resp.Errors.RelationshipsErrors},
			{"theorems", resp.Errors.TheoremsErrors},
		}
		for _, input := range inputs {
			for _, msg := range input.errs {
				fmt.Printf("%v:%v: %v\n", input.name, position(msg.Interval), msg.Message)
			}
		}
		return
	}

	lines := strings.Split(req.Theorems, "\n")
	quote := func(interval *types.BufferInterval) string {
		if interval == nil || interval.Start == nil || interval.Start.Line == 0 || int(interval.Start.Line) > len(lines) {
			return ""
		}
		return strings.TrimSpace(lines[interval.Start.Line-1])
	}
	print := func(result *types.Result, interval *types.BufferInterval) {
		status := "PASS"
		switch result.Status {
		case types.ResultStatus_Reject:
			status = "FAIL"
		case types.ResultStatus_Error:
			status = "ERROR"
		}
		fmt.Printf("%-5v %v: %v", status, position(interval), quote(interval))
		if result.Status == types.ResultStatus_Error && result.Message != "" {
			fmt.Printf(" (%v)", result.Message)
		}
		fmt.Println()
	}
	for _, thm := range resp.Result.AuthorizationTheoremsResult {
		print(thm.Result.Result, thm.Interval)
	}
	for _, thm := range resp.Result.DelegationTheoremsResult {
		print(thm.Result.Result, thm.Interval)
	}
	fmt.Printf("\n%v theorems, %v failed\n", resp.Result.TheoremCount, resp.Result.Failures)
}

func position(interval *types.BufferInterval) string {
	if interval == nil || interval.Start == nil {
		return "?"
	}
	return fmt.Sprintf("%v:%v", interval.Start.Line, interval.Start.Column)
}
