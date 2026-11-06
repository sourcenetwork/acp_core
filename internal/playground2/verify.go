package playground2

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/sourcenetwork/acp_core/internal/alloy"
	pc "github.com/sourcenetwork/acp_core/internal/policy_compiler"
	acperrors "github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// maxParallelChecks bounds the alloy processes running at once
const maxParallelChecks = 4

var ErrVerify = acperrors.New("verify", acperrors.ErrorType_BAD_INPUT)

// alloyErrorHeader matches the first line of an alloy error, eg. `Syntax error in /x/c.als at line 5 column 24:`
var alloyErrorHeader = regexp.MustCompile(`(\w+ error) in \S+ at line (\d+) column (\d+):\s*(.*)$`)

// Verify checks every property of every contract up to the configured scope, each in its own alloy run,
// so an invalid formula leaves its property unverified without affecting the others.
func (g *Analyzer) Verify(ctx context.Context, req *types.VerifyRequest) (*types.VerifyResponse, error) {
	if req.SystemModel == nil {
		return nil, acperrors.Wrap("system model is required", ErrVerify)
	}
	requirements := req.SystemModel.Requirements
	if err := pc.ValidateRequirements(requirements); err != nil {
		return nil, err
	}
	scope := uint32(defaultScope)
	if req.Config != nil && req.Config.Scope != 0 {
		scope = req.Config.Scope
	}
	if scope > maxScope {
		return nil, acperrors.Wrap(fmt.Sprintf("scope %v exceeds the maximum of %v", scope, maxScope), ErrVerify)
	}

	resp := &types.VerifyResponse{Scope: scope}
	type job struct {
		contract *types.Contract
		index    int
		result   *types.PropertyResult
	}
	var jobs []*job
	for _, r := range requirements {
		for _, c := range r.Contracts {
			result := &types.ContractResult{Requirement: r.Name, Contract: c.Name, Status: types.VerifyStatus_DRAFT}
			for i := range c.Properties {
				result.Properties = append(result.Properties, &types.PropertyResult{})
				jobs = append(jobs, &job{contract: c, index: i, result: result.Properties[i]})
			}
			resp.Results = append(resp.Results, result)
		}
	}
	if len(resp.Results) == 0 {
		return nil, acperrors.Wrap("system model has no contracts", ErrVerify)
	}
	if len(jobs) == 0 {
		return resp, nil // only drafts: nothing to compile
	}

	model, err := pc.Compile(req.SystemModel.Policy, types.PolicyMarshalingType_YAML)
	if err != nil {
		return nil, err
	}
	constraints, err := pc.LowerAssumptions(model, req.SystemModel.Constraints)
	if err != nil {
		return nil, err
	}

	dir, err := os.MkdirTemp("", "playground2-verify-")
	if err != nil {
		return nil, fmt.Errorf("verify: %w", err)
	}
	defer os.RemoveAll(dir)
	if err := writeFiles(dir, append(alloy.Emit(model, module), alloy.EmitAssumptions(model, module, constraints))); err != nil {
		return nil, fmt.Errorf("verify: %w", err)
	}

	errs := make([]error, len(jobs))
	sem := make(chan struct{}, maxParallelChecks)
	var wg sync.WaitGroup
	for i, j := range jobs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			errs[i] = g.verifyProperty(ctx, dir, model, j.contract, j.index, scope, j.result)
		}()
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	for _, result := range resp.Results {
		result.Status = contractStatus(result.Properties)
	}
	return resp, nil
}

// contractStatus combines the results of a contract's properties: a violation is definite, so it wins
// over an unverified property. A contract without properties is a draft.
func contractStatus(properties []*types.PropertyResult) types.VerifyStatus {
	if len(properties) == 0 {
		return types.VerifyStatus_DRAFT
	}
	status := types.VerifyStatus_SATISFIED
	for _, p := range properties {
		switch p.Status {
		case types.VerifyStatus_VIOLATED:
			return types.VerifyStatus_VIOLATED
		case types.VerifyStatus_UNVERIFIED:
			status = types.VerifyStatus_UNVERIFIED
		}
	}
	return status
}

// verifyProperty runs the alloy command of a contract's index-th property, filling result.
// Problems with the property leave it unverified, whereas failing to run alloy at all is returned as an error.
func (g *Analyzer) verifyProperty(ctx context.Context, dir string, model *pc.Model, contract *types.Contract, index int, scope uint32, result *types.PropertyResult) error {
	file := alloy.EmitProperty(model, module, contract, index, scope)
	command := alloy.PropertyCommand(contract, index)
	if err := writeFiles(dir, []alloy.File{file}); err != nil {
		return fmt.Errorf("verify: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, g.Timeout)
	defer cancel()
	outDir := filepath.Join(dir, "out", command)
	cmd := exec.CommandContext(ctx, g.AlloyPath, "exec", "-q", "-f", "-c", command, "-t", "xml",
		"-o", outDir, filepath.Join(dir, file.Name))
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		result.Message = fmt.Sprintf("alloy timed out after %v at scope %v: reduce the scope", g.Timeout, scope)
		return nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.Message = formulaError(file, string(out))
		return nil
	}
	if err != nil {
		return acperrors.New(fmt.Sprintf("running alloy: %v", err), acperrors.ErrorType_INTERNAL)
	}

	// alloy writes a solution only if it finds an instance
	data, err := os.ReadFile(filepath.Join(outDir, command+"-solution-0.xml"))
	found := err == nil
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("verify: %w", err)
	}

	isAssert := contract.Properties[index].Type == types.PropertyType_ASSERT
	if found == isAssert {
		result.Status = types.VerifyStatus_VIOLATED
	} else {
		result.Status = types.VerifyStatus_SATISFIED
	}
	if !found {
		return nil
	}

	inst, err := alloy.ParseInstance(model, data)
	if err != nil {
		return fmt.Errorf("verify: %w", err)
	}
	result.World = &types.World{Relationships: strings.Join(inst.Relationships, "\n") + "\n"}
	// skolems are labeled after the command, eg. $contract_doc_edit_1_d
	prefix := "$" + command + "_"
	for label, values := range inst.Skolems {
		result.Bindings = append(result.Bindings, &types.Binding{Name: strings.TrimPrefix(label, prefix), Values: values})
	}
	sort.Slice(result.Bindings, func(i, j int) bool { return result.Bindings[i].Name < result.Bindings[j].Name })
	return nil
}

// formulaError extracts alloy's error message, locating it in the property's formula when possible,
// eg. `syntax error at formula line 1 column 16: The name "nope" cannot be found.`
func formulaError(file alloy.File, out string) string {
	lines := strings.Split(out, "\n")
	for i, line := range lines {
		m := alloyErrorHeader.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		msg := []string{}
		if m[4] != "" {
			msg = append(msg, m[4])
		}
		for _, next := range lines[i+1:] {
			next = strings.TrimSpace(next)
			if next == "" || strings.HasPrefix(next, "at ") || alloyErrorHeader.MatchString(next) {
				break
			}
			msg = append(msg, next)
		}
		where := "outside the formula"
		lineNo, _ := strconv.Atoi(m[2])
		column, _ := strconv.Atoi(m[3])
		if fl, fc, ok := alloy.FormulaPosition(file, lineNo, column); ok {
			where = fmt.Sprintf("at formula line %v column %v", fl, fc)
		}
		return fmt.Sprintf("%v %v: %v", strings.ToLower(m[1]), where, strings.Join(msg, " "))
	}
	return "alloy: " + alloyError(out)
}
