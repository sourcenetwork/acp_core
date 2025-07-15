package services

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/sandbox"
	"github.com/sourcenetwork/acp_core/internal/simulator"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ types.PlaygroundServiceServer = (*playgroundService)(nil)

// playgroundService implements the ACP module MsgServer interface and accepts
// decorating functions which can wrap the execution of a Msg.
type playgroundService struct {
	runtime runtime.RuntimeManager
}

// NewCmdSrever creates a message server for Embedded ACP
func NewPlaygroundService(runtime runtime.RuntimeManager) types.PlaygroundServiceServer {
	return &playgroundService{
		runtime: runtime,
	}
}

func (s *playgroundService) NewSandbox(ctx context.Context, req *types.NewSandboxRequest) (*types.NewSandboxResponse, error) {
	return sandbox.HandleNewSandboxRequest(ctx, s.runtime, req)
}
func (s *playgroundService) ListSandboxes(ctx context.Context, req *types.ListSandboxesRequest) (*types.ListSandboxesResponse, error) {
	return sandbox.HandleListSandboxes(ctx, s.runtime, req)
}
func (s *playgroundService) SetState(ctx context.Context, req *types.SetStateRequest) (*types.SetStateResponse, error) {
	handler := sandbox.SetStateHandler{}
	return handler.Handle(ctx, s.runtime, req)
}
func (s *playgroundService) GetCatalogue(ctx context.Context, req *types.GetCatalogueRequest) (*types.GetCatalogueResponse, error) {
	return sandbox.HandleGetCatalogue(ctx, s.runtime, req)
}
func (s *playgroundService) VerifyTheorems(ctx context.Context, req *types.VerifyTheoremsRequest) (*types.VerifyTheoremsResponse, error) {
	return sandbox.HandleVerifyTheorem(ctx, s.runtime, req)
}

func (s *playgroundService) RestoreScratchpad(ctx context.Context, req *types.RestoreScratchpadRequest) (*types.RestoreScratchpadResponse, error) {
	return sandbox.HandleRestoreScratchpad(ctx, s.runtime, req)
}

func (s *playgroundService) Simulate(ctx context.Context, req *types.SimulateRequest) (*types.SimulateReponse, error) {
	return simulator.HandleSimulate(ctx, s.runtime, req)
}

func (s *playgroundService) GetSandbox(ctx context.Context, req *types.GetSandboxRequest) (*types.GetSandboxResponse, error) {
	return sandbox.HandleGetSandbox(ctx, s.runtime, req)
}

func (s *playgroundService) GetSampleSandboxes(ctx context.Context, req *types.GetSampleSandboxesRequest) (*types.GetSampleSandboxesResponse, error) {
	return sandbox.HandleGetSandboxSamples(ctx, s.runtime, req)
}

func (s *playgroundService) Expand(ctx context.Context, req *types.ExpandRequest) (*types.ExpandResponse, error) {
	return sandbox.HandleExpand(ctx, s.runtime, req)
}
