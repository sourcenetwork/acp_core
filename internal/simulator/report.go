package simulator

import (
	"strings"

	"github.com/sourcenetwork/acp_core/internal/theorem"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func GenerateSimulationReport(result *types.AnnotatedSimulationResult) string {
	if result.Errors.HasErrors() {
		return failureReport(result)
	}
	return sucessReport(result)
}

func failureReport(result *types.AnnotatedSimulationResult) string {
	return result.Errors.String()
}

func sucessReport(result *types.AnnotatedSimulationResult) string {
	builder := strings.Builder{}

	builder.WriteString("policy:\n")
	builder.WriteString(result.Declaration.Policy)
	builder.WriteRune('\n')

	builder.WriteString("relationships:\n")
	builder.WriteString(result.Declaration.RelationshipSet)
	builder.WriteRune('\n')
	builder.WriteRune('\n')

	builder.WriteString("Authorization Theorems:\n")
	for _, thm := range result.PolicyTheoremResult.AuthorizationTheoremsResult {
		if !thm.Result.Result.Valid {
			builder.WriteString("FAIL ")
		}
		builder.WriteRune('\t')
		builder.WriteString(theorem.FormatAuthorizationTheorem(thm.Result.Theorem))
		builder.WriteRune('\n')
	}
	builder.WriteRune('\n')

	builder.WriteString("Delegation Theorems:\n")
	for _, thm := range result.PolicyTheoremResult.DelegationTheoremsResult {
		if !thm.Result.Result.Valid {
			builder.WriteString("FAIL ")
		}
		builder.WriteRune('\t')
		builder.WriteString(theorem.FormatDelegationTheorem(thm.Result.Theorem))
		builder.WriteRune('\n')
	}
	builder.WriteRune('\n')

	builder.WriteString("Reachability Theorems:\n")
	for _, thm := range result.PolicyTheoremResult.ReachabilityTheoremsResult {
		if !thm.Result.Result.Valid {
			builder.WriteString("FAIL ")
		}
		builder.WriteRune('\t')
		builder.WriteString(theorem.FormatReachabilityTheorem(thm.Result.Theorem))
		builder.WriteRune('\n')
	}

	return builder.String()
}
