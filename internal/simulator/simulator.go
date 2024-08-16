package simulator

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/sandbox"
	"github.com/sourcenetwork/acp_core/pkg/playground"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
)

func HandleSimulate(ctx context.Context, manager runtime.RuntimeManager, req *playground.SimulateRequest) (*playground.SimulateReponse, error) {
	manager, err := runtime.NewRuntimeManager(
		runtime.WithLogger(manager.GetLogger()),
	)
	if err != nil {
		return nil, newSimulateError(err)
	}

	newResp, err := sandbox.HandleNewSandboxRequest(ctx, manager, &playground.NewSandboxRequest{})
	if err != nil {
		return nil, newSimulateError(err)
	}

	handler := sandbox.SetStateHandler{}
	setResp, err := handler.Handle(ctx, manager, &playground.SetStateRequest{
		Handle: newResp.Record.Handle,
		Data:   req.Data,
	})
	if err != nil {
		return nil, newSimulateError(err)
	}
	if setResp.Errors.HasErrors() {
		return &playground.SimulateReponse{
			ValidData: false,
			Errors:    setResp.Errors,
		}, nil
	}

	verifyResp, err := sandbox.HandleVerifyTheorem(ctx, manager, &playground.VerifyTheoremsRequest{Handle: newResp.Record.Handle})
	if err != nil {
		return nil, newSimulateError(err)
	}

	return &playground.SimulateReponse{
		ValidData: true,
		Errors:    nil,
		Record:    setResp.Record,
		Result:    verifyResp.Result,
	}, nil
}
