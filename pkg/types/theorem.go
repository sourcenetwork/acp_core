package types

import "github.com/sourcenetwork/acp_core/pkg/utils"

func FromAnnotatedResult(result *AnnotatedPolicyTheoremResult) *PolicyTheoremResult {
	azThms := utils.MapSlice(result.AuthorizationTheoremsResult, func(t *AnnotatedAuthorizationTheoremResult) *AuthorizationTheoremResult {
		return t.Result
	})

	delThms := utils.MapSlice(result.DelegationTheoremsResult, func(t *AnnotatedDelegationTheoremResult) *DelegationTheoremResult {
		return t.Result
	})

	reachThms := utils.MapSlice(result.ReachabilityTheoremsResult, func(t *AnnotatedReachabilityTheoremResult) *ReachabilityTheoremResult {
		return t.Result
	})

	return &PolicyTheoremResult{
		Theorem:                     result.Theorem,
		AuthorizationTheoremsResult: azThms,
		DelegationTheoremsResult:    delThms,
		ReachabilityTheoremsResult:  reachThms,
	}
}
