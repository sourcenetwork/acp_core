package theorem

import (
	"strings"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

func GenerateTheoremReport(result *types.PolicyTheoremResult) string {
	builder := strings.Builder{}

	builder.WriteString("Authorization Theorems:\n")
	for _, thm := range result.AuthorizationTheoremsResult {
		if thm.Result.Status == types.ResultStatus_Reject {
			builder.WriteString("FAIL ")
		} else if thm.Result.Status == types.ResultStatus_Error {
			builder.WriteString("ERROR ")
		}
		builder.WriteRune('\t')
		builder.WriteString(FormatAuthorizationTheorem(thm.Theorem))
		builder.WriteRune('\n')
	}
	builder.WriteRune('\n')

	builder.WriteString("Delegation Theorems:\n")
	for _, thm := range result.DelegationTheoremsResult {
		if thm.Result.Status == types.ResultStatus_Reject {
			builder.WriteString("FAIL ")
		} else if thm.Result.Status == types.ResultStatus_Error {
			builder.WriteString("ERROR ")
		}
		builder.WriteRune('\t')
		builder.WriteString(FormatDelegationTheorem(thm.Theorem))
		builder.WriteRune('\n')
	}
	builder.WriteRune('\n')

	builder.WriteString("Reachability Theorems:\n")
	for _, thm := range result.ReachabilityTheoremsResult {
		if thm.Result.Status == types.ResultStatus_Reject {
			builder.WriteString("FAIL ")
		} else if thm.Result.Status == types.ResultStatus_Error {
			builder.WriteString("ERROR ")
		}
		builder.WriteRune('\t')
		builder.WriteString(FormatReachabilityTheorem(thm.Theorem))
		builder.WriteRune('\n')
	}

	return builder.String()
}
