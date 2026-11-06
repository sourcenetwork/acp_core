package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/sourcenetwork/acp_core/internal/playground2"
	pc "github.com/sourcenetwork/acp_core/internal/policy_compiler"
	"github.com/sourcenetwork/acp_core/pkg/services"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// verify checks contracts against a policy and its constraints, printing a report.
// Counterexamples and examples are printed, and written to <outDir>/<contract>-<property>/relationships.txt if outDir is set.
// It exits with 1 unless every contract is satisfied, drafts aside.
func verify(policyPath, contractsPath, constraintsPath, outDir string, scope uint) {
	requirements, err := pc.UnmarshalRequirements(read(contractsPath))
	if err != nil {
		log.Fatal(err)
	}
	model := &types.SystemModel{Policy: read(policyPath), Requirements: requirements}
	if constraintsPath != "" {
		model.Constraints, err = pc.UnmarshalAssumptions(read(constraintsPath))
		if err != nil {
			log.Fatal(err)
		}
	}

	resp, err := services.NewPlayground2Service("").Verify(context.Background(), &types.VerifyRequest{
		SystemModel: model,
		Config:      &types.VerifyConfig{Scope: uint32(scope)},
	})
	if err != nil {
		log.Fatal(err)
	}

	var save func(string, *types.World)
	if outDir != "" {
		save = func(name string, w *types.World) {
			dir := filepath.Join(outDir, name)
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "relationships.txt"), []byte(w.Relationships), 0644); err != nil {
				log.Fatal(err)
			}
		}
	}
	if playground2.PrintReport(os.Stdout, requirements, resp, save) > 0 {
		os.Exit(1)
	}
}
