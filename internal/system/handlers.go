package system

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func HandleSetParams(ctx context.Context, runtime runtime.RuntimeManager, req *types.SetParamsRequest) (*types.SetParamsResponse, error) {
	_, err := auth.ExtractPrincipalWithType(ctx, auth.Root)
	if err != nil {
		return nil, fmt.Errorf("set params: requires root principal: %w: %v", auth.ErrUnauthorized, err)
	}

	repo := NewParamsRepository(runtime)
	err = repo.Set(ctx, req.Params)
	if err != nil {
		return nil, fmt.Errorf("set params: %v", err)
	}

	return &types.SetParamsResponse{}, nil
}

func HandleGetParams(ctx context.Context, runtime runtime.RuntimeManager, req *types.GetParamsRequest) (*types.GetParamsResponse, error) {
	repo := NewParamsRepository(runtime)
	params, err := repo.GetOrDefault(ctx)
	if err != nil {
		return nil, fmt.Errorf("get params: %v", err)
	}

	return &types.GetParamsResponse{
		Params: params,
	}, nil
}
