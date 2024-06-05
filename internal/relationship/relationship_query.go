package relationship

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func FilterRelationshipsHandler(ctx context.Context, runtime runtime.RuntimeManager, req *types.FilterRelationshipsRequest) (*types.FilterRelationshipsResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	rec, err := engine.GetPolicy(ctx, req.PolicyId)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, fmt.Errorf("policy %v: %w", req.PolicyId, types.ErrPolicyNotFound)
	}

	records, err := engine.FilterRelationships(ctx, rec.Policy, req.Selector)
	if err != nil {
		return nil, err
	}

	return &types.FilterRelationshipsResponse{
		Records: records,
	}, nil
}
