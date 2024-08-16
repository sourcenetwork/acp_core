package types

import (
	"strings"

	"github.com/sourcenetwork/acp_core/pkg/utils"
)

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

// PrettyString formats an AuthorizationThereom back to its DSL representation
func (t *AuthorizationTheorem) PrettyString() string {
	builder := strings.Builder{}
	if !t.AssertTrue {
		builder.WriteRune('!')
	}
	builder.WriteString(t.Operation.Object.Resource)
	builder.WriteRune(':')
	builder.WriteString(t.Operation.Object.Id)
	builder.WriteRune('#')
	builder.WriteString(t.Operation.Permission)
	builder.WriteRune('@')
	builder.WriteString(t.Actor.Id)
	return builder.String()
}

// PrettyString formats an DelegationTheorem back to its DSL representation
func (t *DelegationTheorem) PrettyString() string {
	builder := strings.Builder{}
	if !t.AssertTrue {
		builder.WriteRune('!')
	}
	builder.WriteString(t.Operation.Object.Resource)
	builder.WriteRune(':')
	builder.WriteString(t.Operation.Object.Id)
	builder.WriteRune('#')
	builder.WriteString(t.Operation.Permission)
	builder.WriteString(" > ")
	builder.WriteString(t.Actor.Id)
	return builder.String()
}

// PrettyString formats an ReachabilityTheorem back to its DSL representation
func (t *ReachabilityTheorem) PrettyString() string {
	return ""
}
