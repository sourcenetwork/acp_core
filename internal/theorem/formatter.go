package theorem

import (
	"strings"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

func FormatAuthorizationTheorem(t *types.AuthorizationTheorem) string {
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

func FormatDelegationTheorem(t *types.DelegationTheorem) string {
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

func FormatReachabilityTheorem(t *types.ReachabilityTheorem) string {
	return ""
}
