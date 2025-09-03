package services

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/decorator"
	"github.com/sourcenetwork/acp_core/internal/sandbox"
	"github.com/sourcenetwork/acp_core/internal/simulator"
	"github.com/sourcenetwork/acp_core/internal/telemetry"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ types.PlaygroundServiceServer = (*playgroundService)(nil)

// PlaygroundConfig models configuration params to the playground engine
type PlaygroundConfig struct {
	// PublishErrorEndpoint, if set, will cause all internal errors
	// to be POSTed to endpoint. Assumes the endpoint to be a
	// "PublishError" endpoint as defined in the playground backend service
	PublishErrorEndpoint string
}

// playgroundService implements the ACP module MsgServer interface and accepts
// decorating functions which can wrap the execution of a Msg.
type playgroundService struct {
	runtime runtime.RuntimeManager
	dec     decorator.Decorator
}

// NewCmdSrever creates a message server for Embedded ACP
func NewPlaygroundService(runtime runtime.RuntimeManager, config *PlaygroundConfig) types.PlaygroundServiceServer {
	decorators := []decorator.Decorator{
		decorator.RecoverDecorator,
	}

	if config != nil && config.PublishErrorEndpoint != "" {
		client := telemetry.NewPlaygroundBackendErrorClient(config.PublishErrorEndpoint)
		decorators = append(decorators, sandbox.InternalErrorPublisherDecorator(client))
	}

	// last decorator means it will start off by initializing the request data
	decorators = append(decorators, decorator.RequestDataInitializerDecorator)
	return &playgroundService{
		dec:     decorator.Chain(decorators...),
		runtime: runtime,
	}
}

func (s *playgroundService) NewSandbox(ctx context.Context, req *types.NewSandboxRequest) (*types.NewSandboxResponse, error) {
	h := func(ctx context.Context, req *types.NewSandboxRequest) (*types.NewSandboxResponse, error) {
		return sandbox.HandleNewSandboxRequest(ctx, s.runtime, req)
	}
	return decorator.DecorateTypedHandler(h, s.dec)(ctx, req)
}
func (s *playgroundService) ListSandboxes(ctx context.Context, req *types.ListSandboxesRequest) (*types.ListSandboxesResponse, error) {
	h := func(ctx context.Context, req *types.ListSandboxesRequest) (*types.ListSandboxesResponse, error) {
		return sandbox.HandleListSandboxes(ctx, s.runtime, req)
	}
	return decorator.DecorateTypedHandler(h, s.dec)(ctx, req)
}
func (s *playgroundService) SetState(ctx context.Context, req *types.SetStateRequest) (*types.SetStateResponse, error) {
	h := func(ctx context.Context, req *types.SetStateRequest) (*types.SetStateResponse, error) {
		handler := sandbox.SetStateHandler{}
		return handler.Handle(ctx, s.runtime, req)
	}
	return decorator.DecorateTypedHandler(h, s.dec)(ctx, req)
}
func (s *playgroundService) GetCatalogue(ctx context.Context, req *types.GetCatalogueRequest) (*types.GetCatalogueResponse, error) {
	h := func(ctx context.Context, req *types.GetCatalogueRequest) (*types.GetCatalogueResponse, error) {
		return sandbox.HandleGetCatalogue(ctx, s.runtime, req)
	}
	return decorator.DecorateTypedHandler(h, s.dec)(ctx, req)
}
func (s *playgroundService) VerifyTheorems(ctx context.Context, req *types.VerifyTheoremsRequest) (*types.VerifyTheoremsResponse, error) {
	h := func(ctx context.Context, req *types.VerifyTheoremsRequest) (*types.VerifyTheoremsResponse, error) {
		return sandbox.HandleVerifyTheorem(ctx, s.runtime, req)
	}
	return decorator.DecorateTypedHandler(h, s.dec)(ctx, req)
}

func (s *playgroundService) RestoreScratchpad(ctx context.Context, req *types.RestoreScratchpadRequest) (*types.RestoreScratchpadResponse, error) {
	h := func(ctx context.Context, req *types.RestoreScratchpadRequest) (*types.RestoreScratchpadResponse, error) {
		return sandbox.HandleRestoreScratchpad(ctx, s.runtime, req)
	}
	return decorator.DecorateTypedHandler(h, s.dec)(ctx, req)
}

func (s *playgroundService) Simulate(ctx context.Context, req *types.SimulateRequest) (*types.SimulateResponse, error) {
	h := func(ctx context.Context, req *types.SimulateRequest) (*types.SimulateResponse, error) {
		return simulator.HandleSimulate(ctx, s.runtime, req)
	}
	return decorator.DecorateTypedHandler(h, s.dec)(ctx, req)
}

func (s *playgroundService) GetSandbox(ctx context.Context, req *types.GetSandboxRequest) (*types.GetSandboxResponse, error) {
	h := func(ctx context.Context, req *types.GetSandboxRequest) (*types.GetSandboxResponse, error) {
		return sandbox.HandleGetSandbox(ctx, s.runtime, req)
	}
	return decorator.DecorateTypedHandler(h, s.dec)(ctx, req)
}

func (s *playgroundService) GetSampleSandboxes(ctx context.Context, req *types.GetSampleSandboxesRequest) (*types.GetSampleSandboxesResponse, error) {
	h := func(ctx context.Context, req *types.GetSampleSandboxesRequest) (*types.GetSampleSandboxesResponse, error) {
		return sandbox.HandleGetSandboxSamples(ctx, s.runtime, req)
	}
	return decorator.DecorateTypedHandler(h, s.dec)(ctx, req)
}

func (s *playgroundService) ExplainCheck(ctx context.Context, req *types.ExplainCheckRequest) (*types.ExplainCheckResponse, error) {
	h := func(ctx context.Context, req *types.ExplainCheckRequest) (*types.ExplainCheckResponse, error) {
		return sandbox.HandleExplainCheck(ctx, s.runtime, req)
	}
	return decorator.DecorateTypedHandler(h, s.dec)(ctx, req)
}

func (s *playgroundService) DOTExplainCheck(ctx context.Context, req *types.DOTExplainCheckRequest) (*types.DOTExplainCheckResponse, error) {
	h := func(ctx context.Context, req *types.DOTExplainCheckRequest) (*types.DOTExplainCheckResponse, error) {
		return sandbox.HandleDOTExplainCheck(ctx, s.runtime, req)
	}
	return decorator.DecorateTypedHandler(h, s.dec)(ctx, req)
}
