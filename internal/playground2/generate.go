// package playground2 explores the states of a system model
package playground2

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sourcenetwork/acp_core/internal/alloy"
	pc "github.com/sourcenetwork/acp_core/internal/policy_compiler"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

const (
	defaultCount  = 1
	defaultScope  = 3
	defaultActors = 3
	maxCount      = 50
	maxScope      = 6
	maxActors     = 10

	// module names the compiled system model in the generator's working directory
	module = "system"
	// command names the alloy command which generates worlds
	command = "generate"
)

var ErrGenerate = errors.New("generate", errors.ErrorType_BAD_INPUT)

// Generator generates worlds by running the Alloy analyzer
type Generator struct {
	// AlloyPath is the alloy executable
	AlloyPath string
	// Timeout bounds each Alloy run
	Timeout time.Duration
}

// Generate compiles the system model and asks Alloy for up to count worlds satisfying
// the system constraints, translating each back into relationships.
func (g *Generator) Generate(ctx context.Context, req *types.GenerateRequest) (*types.GenerateResponse, error) {
	if req.SystemModel == nil {
		return nil, errors.Wrap("system model is required", ErrGenerate)
	}
	cfg, err := withDefaults(req.Config)
	if err != nil {
		return nil, err
	}

	model, err := pc.Compile(req.SystemModel.Policy, types.PolicyMarshalingType_YAML)
	if err != nil {
		return nil, err
	}
	constraints, err := pc.LowerAssumptions(model, req.SystemModel.Constraints)
	if err != nil {
		return nil, err
	}

	dir, err := os.MkdirTemp("", "playground2-generate-")
	if err != nil {
		return nil, fmt.Errorf("generate: %w", err)
	}
	defer os.RemoveAll(dir)

	files := append(alloy.Emit(model, module),
		alloy.EmitAssumptions(model, module, constraints),
		alloy.EmitGenerate(model, module, command, cfg.Scope, cfg.Actors),
	)
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f.Name), []byte(f.Content), 0644); err != nil {
			return nil, fmt.Errorf("generate: %w", err)
		}
	}

	outDir := filepath.Join(dir, "out")
	if err := g.run(ctx, dir, outDir, cfg.Count); err != nil {
		return nil, err
	}

	resp := &types.GenerateResponse{}
	for i := 0; i < int(cfg.Count); i++ {
		data, err := os.ReadFile(filepath.Join(outDir, fmt.Sprintf("%v-solution-%v.xml", command, i)))
		if os.IsNotExist(err) {
			break // alloy found no more worlds
		}
		if err != nil {
			return nil, fmt.Errorf("generate: %w", err)
		}
		inst, err := alloy.ParseInstance(model, data)
		if err != nil {
			return nil, fmt.Errorf("generate: %w", err)
		}
		resp.Worlds = append(resp.Worlds, &types.World{
			Relationships: strings.Join(inst.Relationships, "\n") + "\n",
		})
	}
	return resp, nil
}

func (g *Generator) run(ctx context.Context, dir, outDir string, count uint32) error {
	ctx, cancel := context.WithTimeout(ctx, g.Timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, g.AlloyPath, "exec", "-q", "-f", "-c", command, "-t", "xml",
		"-r", fmt.Sprint(count), "-o", outDir, filepath.Join(dir, command+".als"))
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return errors.Wrap(fmt.Sprintf("alloy timed out after %v: reduce the scope or the number of actors", g.Timeout), ErrGenerate)
	}
	if err != nil {
		return errors.New(fmt.Sprintf("running alloy: %v: %v", err, alloyError(string(out))), errors.ErrorType_INTERNAL)
	}
	return nil
}

// alloyError extracts the relevant lines of alloy's output, dropping stack traces
func alloyError(out string) string {
	var lines []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "at ") {
			continue
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, " ")
}

func withDefaults(cfg *types.GenerateConfig) (*types.GenerateConfig, error) {
	out := &types.GenerateConfig{Count: defaultCount, Scope: defaultScope, Actors: defaultActors}
	if cfg != nil {
		if cfg.Count != 0 {
			out.Count = cfg.Count
		}
		if cfg.Scope != 0 {
			out.Scope = cfg.Scope
		}
		if cfg.Actors != 0 {
			out.Actors = cfg.Actors
		}
	}
	if out.Count > maxCount {
		return nil, errors.Wrap(fmt.Sprintf("count %v exceeds the maximum of %v", out.Count, maxCount), ErrGenerate)
	}
	if out.Scope > maxScope {
		return nil, errors.Wrap(fmt.Sprintf("scope %v exceeds the maximum of %v", out.Scope, maxScope), ErrGenerate)
	}
	if out.Actors > maxActors {
		return nil, errors.Wrap(fmt.Sprintf("actors %v exceeds the maximum of %v", out.Actors, maxActors), ErrGenerate)
	}
	return out, nil
}
