package playground

import (
	"strings"
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/playground"
)

type Assertion func(*testing.T, *playground.SandboxDataErrors)

func HasPolicyError(msg string) Assertion {
	return func(t *testing.T, errs *playground.SandboxDataErrors) {
		AssertContainsErrWithMsg(t, msg, errs.PolicyErrors)
	}
}

func HasRelationshipsError(msg string) Assertion {
	return func(t *testing.T, errs *playground.SandboxDataErrors) {
		AssertContainsErrWithMsg(t, msg, errs.RelationshipsErrors)
	}
}

func HasTheoremError(msg string) Assertion {
	return func(t *testing.T, errs *playground.SandboxDataErrors) {
		AssertContainsErrWithMsg(t, msg, errs.TheoremsErrors)
	}
}

func AssertContainsErrWithMsg(t *testing.T, want string, msgs []*errors.ParserMessage) {
	var had []string
	for _, msg := range msgs {
		had = append(had, msg.Message)
		if strings.Contains(msg.Message, want) {
			return
		}
	}
	t.Errorf("error not found in parser messages: want %v had %v", want, had)
}
