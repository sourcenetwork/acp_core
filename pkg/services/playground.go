package services

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/sandbox"
	"github.com/sourcenetwork/acp_core/pkg/playground"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
)

var _ playground.PlaygroundServiceServer = (*playgroundService)(nil)

// playgroundService implements the ACP module MsgServer interface and accepts
// decorating functions which can wrap the execution of a Msg.
type playgroundService struct {
	runtime runtime.RuntimeManager
}

// NewCmdSrever creates a message server for Embedded ACP
func NewPlaygroundService(runtime runtime.RuntimeManager) playground.PlaygroundServiceServer {
	return &playgroundService{
		runtime: runtime,
	}
}

func (s *playgroundService) NewSandbox(ctx context.Context, req *playground.NewSandboxRequest) (*playground.NewSandboxResponse, error) {
	return sandbox.HandleNewSandboxRequest(ctx, s.runtime, req)
}
func (s *playgroundService) ListSandboxes(ctx context.Context, req *playground.ListSandboxesRequest) (*playground.ListSandboxesResponse, error) {
	return sandbox.HandleListSandboxes(ctx, s.runtime, req)
}
func (s *playgroundService) SetState(ctx context.Context, req *playground.SetStateRequest) (*playground.SetStateResponse, error) {
	handler := sandbox.SetStateHandler{}
	return handler.Handle(ctx, s.runtime, req)
}
func (s *playgroundService) GetCatalogue(ctx context.Context, req *playground.GetCatalogueRequest) (*playground.GetCatalogueResponse, error) {
	return nil, nil
}
func (s *playgroundService) VerifyTheorems(ctx context.Context, req *playground.VerifyTheoremsRequest) (*playground.VerifyTheoremsResponse, error) {
	return sandbox.HandleVerifyTheorem(ctx, s.runtime, req)
}
