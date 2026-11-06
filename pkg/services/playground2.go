package services

import (
	"context"
	"time"

	"github.com/sourcenetwork/acp_core/internal/decorator"
	"github.com/sourcenetwork/acp_core/internal/playground2"
	"github.com/sourcenetwork/acp_core/internal/telemetry"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ types.Playground2Server = (*Playground2Service)(nil)

const defaultAlloyTimeout = 30 * time.Second

// Playground2Service explores the states of a system model.
// Generate and Verify require the alloy executable; Check and Explain don't.
type Playground2Service struct {
	decorator decorator.Decorator
	analyzer  *playground2.Analyzer
	compiler  *PolicyCompilerService
}

// NewPlayground2Service returns a Playground2Service running alloy from alloyPath.
// An empty alloyPath looks up alloy in PATH.
func NewPlayground2Service(alloyPath string) *Playground2Service {
	if alloyPath == "" {
		alloyPath = "alloy"
	}
	return &Playground2Service{
		decorator: decorator.Chain(telemetry.RequestDataInitializerDecorator, decorator.RecoverDecorator),
		analyzer:  &playground2.Analyzer{AlloyPath: alloyPath, Timeout: defaultAlloyTimeout},
		compiler:  NewPolicyCompilerService(),
	}
}

// Generate returns worlds which are valid states of a system model
func (s *Playground2Service) Generate(ctx context.Context, req *types.GenerateRequest) (*types.GenerateResponse, error) {
	return decorator.DecorateTypedHandler(s.analyzer.Generate, s.decorator)(ctx, req)
}

// Check returns whether an actor has a permission or relation on an object, in a world
func (s *Playground2Service) Check(ctx context.Context, req *types.CheckRequest) (*types.CheckResponse, error) {
	return decorator.DecorateTypedHandler(playground2.Check, s.decorator)(ctx, req)
}

// Explain performs a Check, returning the engine's search graph alongside the result
func (s *Playground2Service) Explain(ctx context.Context, req *types.ExplainRequest) (*types.ExplainResponse, error) {
	return decorator.DecorateTypedHandler(playground2.Explain, s.decorator)(ctx, req)
}

// Verify checks the contracts of a system model against every valid state, up to a scope
func (s *Playground2Service) Verify(ctx context.Context, req *types.VerifyRequest) (*types.VerifyResponse, error) {
	return decorator.DecorateTypedHandler(s.analyzer.Verify, s.decorator)(ctx, req)
}

// Compile returns the Alloy modules a policy compiles to
func (s *Playground2Service) Compile(ctx context.Context, req *types.CompileRequest) (*types.CompileResponse, error) {
	return s.compiler.Compile(ctx, req)
}
