package policy

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
	"github.com/stretchr/testify/require"
)

func TestValidatePolicy_ValidPolicyOk(t *testing.T) {
	ctx := test.NewTestCtx(t)

	pol := `
name: test
resources:
  foo:
    relations:
      reader:
    permissions:
      read:
        expr: reader
`
	resp, err := ctx.Engine.ValidatePolicy(ctx, &types.ValidatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	})
	require.NoError(t, err)
	want := &types.ValidatePolicyResponse{
		ErrorMsg: "",
		Valid:    true,
	}
	require.Equal(t, want, resp)
}

func TestValidatePolicy_InvalidPolicyReturnsErrorMsg(t *testing.T) {
	ctx := test.NewTestCtx(t)

	pol := `
name: test
spec: defra
resources:
  foo:
    relations:
      reader:
    permissions:
      read:
        expr: reader
`
	resp, err := ctx.Engine.ValidatePolicy(ctx, &types.ValidatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	})
	require.NoError(t, err)
	require.False(t, resp.Valid)
	require.Contains(t, resp.ErrorMsg, "defra policy specification")
}
