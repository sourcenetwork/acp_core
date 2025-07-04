package theorem

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/authorizer"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/parser/theorem_parser"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

// NewEvaluator returns an Evaluator bound to a Zanzi instance
func NewEvaluator(zanzi *zanzi.Adapter) *Evaluator {
	return &Evaluator{
		zanzi: zanzi,
	}
}

// Evaluator verifies the correctness of PolicyTheorems
type Evaluator struct {
	zanzi *zanzi.Adapter
}

func (e *Evaluator) evaluateAuthorizationTheorem(ctx context.Context, policy *types.Policy, theorem *types.AuthorizationTheorem) (*types.AuthorizationTheoremResult, error) {
	ok, err := e.zanzi.Check(ctx, policy, theorem.Operation, theorem.Actor)
	if err != nil {
		// FIXME once Zanzi errors are sorted out, assert that the error isn't an IO or internal error
		return &types.AuthorizationTheoremResult{
			Theorem: theorem,
			Result: &types.Result{
				Status:  types.ResultStatus_Error,
				Message: err.Error(),
			},
		}, nil
	}
	return &types.AuthorizationTheoremResult{
		Theorem: theorem,
		Result: &types.Result{
			Status:  toStatus(nxor(ok, theorem.AssertTrue)),
			Message: "",
		},
	}, nil
}

func (e *Evaluator) evaluateReacheabilityTheorem(ctx context.Context, polId *types.Policy, theorem *types.ReachabilityTheorem) (*types.ReachabilityTheoremResult, error) {
	return &types.ReachabilityTheoremResult{
		Result: &types.Result{
			Status:  types.ResultStatus_Accept,
			Message: "",
		},
		Theorem: theorem,
	}, nil
}

func (e *Evaluator) evalDelegationTheorem(ctx context.Context, polId *types.Policy, theorem *types.DelegationTheorem) (*types.DelegationTheoremResult, error) {
	req := authorizer.ManagementRequest{
		Policy:   polId,
		Object:   theorem.Operation.Object,
		Relation: theorem.Operation.Permission,
		Actor:    theorem.Actor,
	}
	authz := authorizer.NewOperationAuthorizer(e.zanzi)
	authorized, err := authz.IsAuthorized(ctx, &req)
	if err != nil {
		// if error is not internal, then user might've supplied invalid data
		// which shouldn't cause the whole execution to fail
		if acpErr, ok := err.(*errors.Error); ok && acpErr.Kind != errors.ErrorType_INTERNAL {
			return &types.DelegationTheoremResult{
				Result: &types.Result{
					Status:  types.ResultStatus_Error,
					Message: acpErr.Error(),
				},
				Theorem: theorem,
			}, nil
		}
	}
	return &types.DelegationTheoremResult{
		Result: &types.Result{
			Status:  toStatus(nxor(authorized, theorem.AssertTrue)),
			Message: "",
		},
		Theorem: theorem,
	}, nil
}

// EvalutePolicyTheoremDSL evaluates a PolicyTheorem represented as a string,
// and return a Report which references the text location of each theorem alongside its result
func (e *Evaluator) EvaluatePolicyTheoremDSL(ctx context.Context, polId string, theoremDSL string) (*types.AnnotatedPolicyTheoremResult, error) {
	indexedTheorem, report := theorem_parser.ParsePolicyTheorem(theoremDSL)
	if report.HasError() {
		return nil, report
	}
	policyResult, err := e.EvaluatePolicyTheorem(ctx, polId, indexedTheorem.ToPolicyTheorem())
	if err != nil {
		return nil, err
	}
	annotatedResult := &types.AnnotatedPolicyTheoremResult{
		Theorem: indexedTheorem.ToPolicyTheorem(),
	}
	for i, theorem := range indexedTheorem.AuthorizationTheorems {
		result := policyResult.AuthorizationTheoremsResult[i]
		annotatedAuth := &types.AnnotatedAuthorizationTheoremResult{
			Result:   result,
			Interval: theorem.Interval,
		}
		annotatedResult.AuthorizationTheoremsResult = append(annotatedResult.AuthorizationTheoremsResult, annotatedAuth)
	}
	for i, theorem := range indexedTheorem.DelegationTheorems {
		result := policyResult.DelegationTheoremsResult[i]
		thmResult := &types.AnnotatedDelegationTheoremResult{
			Result:   result,
			Interval: theorem.Interval,
		}
		annotatedResult.DelegationTheoremsResult = append(annotatedResult.DelegationTheoremsResult, thmResult)
	}
	for i, theorem := range indexedTheorem.ReachabilityTheorems {
		result := policyResult.ReachabilityTheoremsResult[i]
		thmResult := &types.AnnotatedReachabilityTheoremResult{
			Result:   result,
			Interval: theorem.Interval,
		}
		annotatedResult.ReachabilityTheoremsResult = append(annotatedResult.ReachabilityTheoremsResult, thmResult)
	}
	return annotatedResult, nil
}

// EvaluatePolicyThereom verifies the whether a PolicyTheorem is correct.
func (e *Evaluator) EvaluatePolicyTheorem(ctx context.Context, polId string, theorem *types.PolicyTheorem) (*types.PolicyTheoremResult, error) {
	rec, err := e.zanzi.GetPolicy(ctx, polId)
	if err != nil {
		return nil, newEvaluatorErr(err)
	}
	if rec == nil {
		return nil, newEvaluatorErr(errors.ErrPolicyNotFound(polId))
	}

	authzResults, err := utils.MapFailableSlice(theorem.AuthorizationTheorems, func(thm *types.AuthorizationTheorem) (*types.AuthorizationTheoremResult, error) {
		return e.evaluateAuthorizationTheorem(ctx, rec.Policy, thm)
	})
	if err != nil {
		return nil, newEvaluatorErr(err)
	}

	delegationResults, err := utils.MapFailableSlice(theorem.DelegationTheorems, func(thm *types.DelegationTheorem) (*types.DelegationTheoremResult, error) {
		return e.evalDelegationTheorem(ctx, rec.Policy, thm)
	})
	if err != nil {
		return nil, newEvaluatorErr(err)
	}

	reachabilityResults, err := utils.MapFailableSlice(theorem.ReachabilityTheorems, func(thm *types.ReachabilityTheorem) (*types.ReachabilityTheoremResult, error) {
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

// nxor implements a not xor function
func nxor(a, b bool) bool {
	return a && b || !a && !b
}

func toStatus(success bool) types.ResultStatus {
	if success {
		return types.ResultStatus_Accept
	}
	return types.ResultStatus_Reject
}
