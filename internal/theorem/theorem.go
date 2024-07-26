package theorem

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/relationship"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

func NewEvaluator(zanzi *zanzi.Adapter) *Evaluator {
	return &Evaluator{
		zanzi: zanzi,
	}
}

type Evaluator struct {
	zanzi *zanzi.Adapter
}

func (e *Evaluator) evaluateAuthorizationTheorem(ctx context.Context, policy *types.Policy, theorem *types.AuthorizationTheorem) (*types.Result, error) {
	ok, err := e.zanzi.Check(ctx, policy, theorem.Operation, theorem.Actor)
	if err != nil {
		// FIXME once Zanzi errors are sorted out, assert that the error isn't an IO or internal error
		return &types.Result{
			Valid:   false,
			Message: err.Error(),
		}, nil
	}
	return &types.Result{
		Valid:   ok,
		Message: "",
	}, nil
}

func (e *Evaluator) evaluateReacheabilityTheorem(ctx context.Context, polId *types.Policy, theorem *types.ReachabilityTheorem) (*types.Result, error) {
	return &types.Result{
		Valid:   true,
		Message: "",
	}, nil
}

func (e *Evaluator) evalDelegationTheorem(ctx context.Context, polId *types.Policy, theorem *types.DelegationTheorem) (*types.Result, error) {
	authorizer := relationship.NewRelationshipAuthorizer(e.zanzi)
	authorized, err := authorizer.IsAuthorized(ctx, polId, theorem.Operation, theorem.Actor)
	if err != nil {
		// if error is not internal, then user might've supplied invalid data
		// which shouldn't cause the whole execution to fail
		if acpErr, ok := err.(*errors.Error); ok && acpErr.Type() != errors.ErrorType_INTERNAL {
			return &types.Result{
				Valid:   false,
				Message: acpErr.Error(),
			}, nil
		}
		return nil, err
	}
	return &types.Result{
		Valid:   authorized,
		Message: "",
	}, nil
}

func (e *Evaluator) EvaluatePolicyTheorem(ctx context.Context, polId string, theorem *types.PolicyTheorem) (*types.PolicyTheoremResult, error) {
	rec, err := e.zanzi.GetPolicy(ctx, polId)
	if err != nil {
		return nil, newEvaluatorErr(err)
	}
	if rec == nil {
		return nil, newEvaluatorErr(errors.NewPolicyNotFound(polId))
	}

	authzResults, err := utils.MapFailableSlice(theorem.AuthorizationThereoms, func(thm *types.AuthorizationTheorem) (*types.Result, error) {
		return e.evaluateAuthorizationTheorem(ctx, rec.Policy, thm)
	})
	if err != nil {
		return nil, newEvaluatorErr(err)
	}

	delegationResults, err := utils.MapFailableSlice(theorem.DelegationTheorems, func(thm *types.DelegationTheorem) (*types.Result, error) {
		return e.evalDelegationTheorem(ctx, rec.Policy, thm)
	})
	if err != nil {
		return nil, newEvaluatorErr(err)
	}

	reachabilityResults, err := utils.MapFailableSlice(theorem.ReachabilityTheorems, func(thm *types.ReachabilityTheorem) (*types.Result, error) {
		return e.evaluateReacheabilityTheorem(ctx, rec.Policy, thm)
	})
	if err != nil {
		return nil, newEvaluatorErr(err)
	}

	return &types.PolicyTheoremResult{
		ReachabilityTheoremsResult:  reachabilityResults,
		DelegationTheoremsResult:    delegationResults,
		AuthorizationTheoremsResult: authzResults,
	}, nil
}

// TheoremGenerator receives a Policy and produces all valid theorems for that Policy
type TheoremGenerator interface {
	GenAuthTheorems(ctx context.Context, polId string) ([]*types.AuthorizationTheorem, error)
	GenAdminTheorems(ctx context.Context, polId string) ([]*types.DelegationTheorem, error)
	GenReachabilityTheorems(ctx context.Context, polId string) ([]*types.ReachabilityTheorem, error)
	GenPolicyTheorem(ctx context.Context, polId string) (*types.PolicyTheorem, error)
}
