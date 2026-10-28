package simulator

import (
	"context"
	"strings"

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

// emptyTheorem is evaluated when a request has no theorems
const emptyTheorem = "Authorizations {}\nDelegations {}"

// HandleEvaluate evaluates a policy, relationships and theorems in an ephemeral sandbox
func HandleEvaluate(ctx context.Context, manager runtime.RuntimeManager, req *types.EvaluateRequest) (*types.EvaluateResponse, error) {
	theorems := req.Theorems
	if strings.TrimSpace(theorems) == "" {
		theorems = emptyTheorem
	}
	resp, err := HandleSimulate(ctx, manager, &types.SimulateRequest{
		Data: &types.SandboxData{
			PolicyDefinition: req.Policy,
			Relationships:    req.Relationships,
			PolicyTheorem:    theorems,
		},
	})
	if err != nil {
		return nil, err
	}
	if !resp.ValidData {
		return &types.EvaluateResponse{Errors: resp.Errors}, nil
	}
	return &types.EvaluateResponse{
		Valid:  true,
		Ctx:    resp.Record.Ctx,
		Result: resp.Result,
	}, nil
}
