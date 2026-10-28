package services

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/decorator"
	"github.com/sourcenetwork/acp_core/internal/simulator"
	"github.com/sourcenetwork/acp_core/internal/telemetry"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ types.SimulatorServer = (*SimulatorService)(nil)

// SimulatorService evaluates policies, relationships and theorems statelessly
type SimulatorService struct {
	decorator decorator.Decorator
	runtime   runtime.RuntimeManager
}

func NewSimulatorService(runtime runtime.RuntimeManager) *SimulatorService {
	return &SimulatorService{
		decorator: decorator.Chain(telemetry.RequestDataInitializerDecorator, decorator.RecoverDecorator),
		runtime:   runtime,
	}
}

// Evaluate loads a policy and relationships in an ephemeral environment and evaluates the theorems against them
func (s *SimulatorService) Evaluate(ctx context.Context, req *types.EvaluateRequest) (*types.EvaluateResponse, error) {
	h := func(ctx context.Context, req *types.EvaluateRequest) (*types.EvaluateResponse, error) {
		return simulator.HandleEvaluate(ctx, s.runtime, req)
	}
	return decorator.DecorateTypedHandler(h, s.decorator)(ctx, req)
}
