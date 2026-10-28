package services

import (
	"context"
	"time"

	"github.com/sourcenetwork/acp_core/internal/decorator"
	"github.com/sourcenetwork/acp_core/internal/playground2"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ types.Playground2Server = (*Playground2Service)(nil)

const defaultAlloyTimeout = 30 * time.Second

// Playground2Service explores the states of a system model.
// Generating worlds requires the alloy executable.
type Playground2Service struct {
	decorator decorator.Decorator
	generator *playground2.Generator
}

// NewPlayground2Service returns a Playground2Service running alloy from alloyPath.
// An empty alloyPath looks up alloy in PATH.
func NewPlayground2Service(alloyPath string) *Playground2Service {
	if alloyPath == "" {
		alloyPath = "alloy"
	}
	return &Playground2Service{
		decorator: decorator.RecoverDecorator,
		generator: &playground2.Generator{AlloyPath: alloyPath, Timeout: defaultAlloyTimeout},
	}
}

// Generate returns worlds which are valid states of a system model
func (s *Playground2Service) Generate(ctx context.Context, req *types.GenerateRequest) (*types.GenerateResponse, error) {
	return decorator.DecorateTypedHandler(s.generator.Generate, s.decorator)(ctx, req)
}
