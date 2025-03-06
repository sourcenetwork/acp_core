package policy

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/theorem"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func HandleGetPolicy(ctx context.Context, runtime runtime.RuntimeManager, req *types.GetPolicyRequest) (*types.GetPolicyResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	rec, err := engine.GetPolicy(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, errors.ErrPolicyNotFound(req.Id)
	}

	return &types.GetPolicyResponse{
		Record: rec,
	}, nil
}

func ListPolicies(ctx context.Context, runtime runtime.RuntimeManager, req *types.ListPoliciesRequest) (*types.ListPoliciesResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	records, err := engine.ListPolicies(ctx)
	if err != nil {
		return nil, err
	}

	return &types.ListPoliciesResponse{
		Records: records,
	}, nil
}

func ValidatePolicy(ctx context.Context, runtime runtime.RuntimeManager, req *types.ValidatePolicyRequest) (*types.ValidatePolicyResponse, error) {
	resp := &types.ValidatePolicyResponse{
		Valid: false,
	}

	ir, err := Unmarshal(req.Policy, req.MarshalType)
	if err != nil {
		resp.ErrorMsg = err.Error()
		return resp, nil
	}

	pol, err := mapIRIntoPolicy(ir, 0)
	if err != nil {
		return nil, fmt.Errorf("ValidatePolicy: %v", err)
	}

	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, fmt.Errorf("ValidatePolicy: %v", err)
	}

	valid, msg, err := engine.ValidatePolicy(ctx, pol)
	if err != nil {
		return nil, fmt.Errorf("ValidatePolicy: %v", err)
	}
	resp.Valid = valid
	resp.ErrorMsg = msg

	return resp, nil
}

func EvaluateTheorem(ctx context.Context, manager runtime.RuntimeManager, req *types.EvaluateTheoremRequest) (*types.EvaluateTheoremResponse, error) {
	engine, err := zanzi.NewZanzi(manager.GetKVStore(), manager.GetLogger())
	if err != nil {
		return nil, newEvaluateTheoremErr(err)
	}

	evaluator := theorem.NewEvaluator(engine)
	result, err := evaluator.EvaluatePolicyTheoremDSL(ctx, req.PolicyId, req.PolicyTheorem)
	if err != nil {
		return nil, newEvaluateTheoremErr(err)
	}

	return &types.EvaluateTheoremResponse{
		Result: result,
	}, nil
}

func GetPolicyCatalogue(ctx context.Context, manager runtime.RuntimeManager, req *types.GetPolicyCatalogueRequest) (*types.GetPolicyCatalogueResponse, error) {
	engine, err := zanzi.NewZanzi(manager.GetKVStore(), manager.GetLogger())
	if err != nil {
		return nil, newEvaluateTheoremErr(err)
	}

	catalogue, err := BuildCatalogue(ctx, engine, req.PolicyId)
	if err != nil {
		return nil, newPolicyCatalogueErr(err)
	}

	return &types.GetPolicyCatalogueResponse{
		Catalogue: catalogue,
	}, nil
}
