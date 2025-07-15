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
		return nil, errors.Wrap("verifying access request", errors.ErrPolicyNotFound(req.PolicyId))
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

func Expand(ctx context.Context, runtime runtime.RuntimeManager, policyId string, op *types.Operation) (string, error) {
	zanzi, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return "", err
	}

	rec, err := zanzi.GetPolicy(ctx, policyId)
	if err != nil {
		return "", err
	}
	if rec == nil {
		return "", errors.Wrap("expand", errors.ErrPolicyNotFound(policyId))
	}

	if op == nil {
		return "", errors.Wrap("expand: invalid operation", errors.ErrorType_BAD_INPUT)
	}
	if op.Permission == "" {
		return "", errors.Wrap("expand: invalid operation permission", errors.ErrorType_BAD_INPUT)
	}
	if op.Object == nil || op.Object.Id == "" || op.Object.Resource == "" {
		return "", errors.Wrap("expand: invalid object", errors.ErrorType_BAD_INPUT)
	}

	tree, err := zanzi.Expand(ctx, rec.Policy, op.Object, op.Permission)
	if err != nil {
		return "", errors.Wrap("expand failed", err)
	}
	return tree, nil
}
