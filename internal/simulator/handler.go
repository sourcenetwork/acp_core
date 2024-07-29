package simulator

import (
	"context"

	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func HandleSimulateRequest(ctx context.Context, req *types.SimulateRequest) (*types.SimulateResponse, error) {
	manager, err := runtime.NewRuntimeManager(runtime.WithMemKV())
	if err != nil {
		return nil, newSimulateErr(err)
	}
	result, err := SimulateDeclaration(ctx, manager, req.Declaration)
	if err != nil {
		return nil, newSimulateErr(err)
	}
	return &types.SimulateResponse{
		Result: result,
	}, nil
}
