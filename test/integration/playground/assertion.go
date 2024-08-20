package playground

import (
	"strings"
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

type Assertion func(*testing.T, *types.SandboxDataErrors)

func HasPolicyError(msg string) Assertion {
	return func(t *testing.T, errs *types.SandboxDataErrors) {
		AssertContainsErrWithMsg(t, msg, errs.PolicyErrors)
	}
}

func HasRelationshipsError(msg string) Assertion {
	return func(t *testing.T, errs *types.SandboxDataErrors) {
		AssertContainsErrWithMsg(t, msg, errs.RelationshipsErrors)
	}
}

func HasTheoremError(msg string) Assertion {
	return func(t *testing.T, errs *types.SandboxDataErrors) {
		AssertContainsErrWithMsg(t, msg, errs.TheoremsErrors)
	}
}

func AssertContainsErrWithMsg(t *testing.T, want string, msgs []*types.LocatedMessage) {
	var had []string
	for _, msg := range msgs {
		had = append(had, msg.Message)
		if strings.Contains(msg.Message, want) {
			return
		}
	}
	t.Errorf("error not found in parser messages: want %v had %v", want, had)
}
