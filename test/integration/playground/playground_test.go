package playground

import (
	"testing"

	_ "github.com/stretchr/testify/require"

	_ "github.com/sourcenetwork/acp_core/pkg/errors"
	_ "github.com/sourcenetwork/acp_core/pkg/types"
	_ "github.com/sourcenetwork/acp_core/test"
)

func Test_NewSandbox_ReturnsHandle(t *testing.T) {
}

func Test_NewSandbox_UnamedSandbox_ReturnsSandboxWithHandleAsName(t *testing.T) {
}

func Test_NewSandbox_CanCreateSandboxWithoutDescription(t *testing.T) {
}

func Test_SetState_EmptyTheoremIsAccepted(t *testing.T) {
}

func Test_Evaluate_SandboxWithEmptyTheoremOk(t *testing.T) {
}

func Test_Evaluate_UninitializedSandboxCannotBeEvaluated(t *testing.T) {
}
