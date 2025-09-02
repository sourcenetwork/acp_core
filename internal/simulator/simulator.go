package simulator

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/sandbox"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func HandleSimulate(ctx context.Context, manager runtime.RuntimeManager, req *types.SimulateRequest) (*types.SimulateResponse, error) {
	manager, err := runtime.NewRuntimeManager(
		runtime.WithLogger(manager.GetLogger()),
	)
	if err != nil {
		return nil, newSimulateError(err)
	}

	newResp, err := sandbox.HandleNewSandboxRequest(ctx, manager, &types.NewSandboxRequest{})
	if err != nil {
		return nil, newSimulateError(err)
	}

	handler := sandbox.SetStateHandler{}
	setResp, err := handler.Handle(ctx, manager, &types.SetStateRequest{
		Handle: newResp.Record.Handle,
		Data:   req.Data,
	})
	if err != nil {
		return nil, newSimulateError(err)
	}
	if setResp.Errors.HasErrors() {
		return &types.SimulateResponse{
			ValidData: false,
			Errors:    setResp.Errors,
		}, nil
	}

	verifyResp, err := sandbox.HandleVerifyTheorem(ctx, manager, &types.VerifyTheoremsRequest{Handle: newResp.Record.Handle})
	if err != nil {
		return nil, newSimulateError(err)
	}

	return &types.SimulateResponse{
		ValidData: true,
		Errors:    nil,
		Record:    setResp.Record,
		Result:    verifyResp.Result,
	}, nil
}
