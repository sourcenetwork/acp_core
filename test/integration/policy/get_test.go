package policy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func TestQueryPolicy_UnknownPolicyReturnsPolicyNotFoundErr(t *testing.T) {
	ctx := test.NewTestCtx(t)

	req := types.GetPolicyRequest{
		Id: "blahblabh",
	}

	resp, err := ctx.Engine.GetPolicy(ctx, &req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_NOT_FOUND)
}

func TestGetPolicy_ReturnsAnExistingPolicy(t *testing.T) {
	ctx := test.NewTestCtx(t)

	pol := `
name: policy
`
	action := test.CreatePolicyAction{
		Policy: pol,
	}
	policy := action.Run(ctx)

	req := types.GetPolicyRequest{
		Id: policy.Id,
	}
	resp, err := ctx.Engine.GetPolicy(ctx, &req)

	want := &types.GetPolicyResponse{
		Policy: &types.Policy{
			Id:           "bc7eb5a8c500111b2459a92ae23f4848537e49599df1b8d70636b5aacb47bd5f",
			Name:         "policy",
			CreationTime: test.DefaultTs,
			ActorResource: &types.ActorResource{
				Name: "actor",
			},
		},
		PolicyRaw:   pol,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	}
	require.Equal(t, want, resp)
	require.NoError(t, err)
}
