package authz_db

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func VerifyAccessRequest(ctx context.Context, runtime runtime.RuntimeManager, req *types.VerifyAccessRequestRequest) (*types.VerifyAccessRequestResponse, error) {
	zanzi, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	rec, err := zanzi.GetPolicy(ctx, req.PolicyId)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, errors.Wrap("verifying access request", errors.NewPolicyNotFound(req.PolicyId))
	}

	for _, op := range req.AccessRequest.Operations {
		ok, err := zanzi.Check(ctx, rec.Policy, op, req.AccessRequest.Actor)
		if err != nil {
			return nil, errors.Wrap("verify access request", err)
		}
		if !ok {
			return &types.VerifyAccessRequestResponse{
				Valid: false,
			}, nil
		}
	}

	return &types.VerifyAccessRequestResponse{
		Valid: true,
	}, nil
}
