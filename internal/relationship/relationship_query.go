package relationship

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func FilterRelationshipsHandler(ctx context.Context, runtime runtime.RuntimeManager, req *types.FilterRelationshipsRequest) (*types.FilterRelationshipsResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, newFilterRelationshpErr(err)
	}

	rec, err := engine.GetPolicy(ctx, req.PolicyId)
	if err != nil {
		return nil, newFilterRelationshpErr(err)
	}
	if rec == nil {
		return nil, newFilterRelationshpErr(errors.ErrPolicyNotFound(req.PolicyId))
	}

	records, err := engine.FilterRelationships(ctx, rec.Policy, req.Selector)
	if err != nil {
		return nil, newFilterRelationshpErr(err)
	}

	return &types.FilterRelationshipsResponse{
		Records: records,
	}, nil
}

func ValidateRelationship(ctx context.Context, manager runtime.RuntimeManager, policyId string, relationship *types.Relationship) (valid bool, errorMsg string, err error) {
	engine, err := zanzi.NewZanzi(manager.GetKVStore(), manager.GetLogger())
	if err != nil {
		return false, "", newValidateRelationshipErr(err)
	}

	rec, err := engine.GetPolicy(ctx, policyId)
	if err != nil {
		return false, "", newValidateRelationshipErr(err)
	}
	if rec == nil {
		return false, "", newValidateRelationshipErr(errors.ErrPolicyNotFound(policyId))
	}

	err = relationshipSpec(rec.Policy, relationship)
	if err != nil {
		return false, err.Error(), nil
	}

	valid, errMsg, err := engine.ValidateRelationship(ctx, rec.Policy, relationship)
	if err != nil {
		return false, "", newValidateRelationshipErr(errors.ErrPolicyNotFound(policyId))
	}

	return valid, errMsg, nil
}
