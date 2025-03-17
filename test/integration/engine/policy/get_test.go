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
spec: none
`
	action := test.CreatePolicyAction{
		Policy: pol,
	}
	policy := action.Run(ctx)

	req := types.GetPolicyRequest{
		Id: policy.Id,
	}
	resp, err := ctx.Engine.GetPolicy(ctx, &req)

	wantPolicy := &types.Policy{
		Id:                "bc7eb5a8c500111b2459a92ae23f4848537e49599df1b8d70636b5aacb47bd5f",
		Name:              "policy",
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		ActorResource: &types.ActorResource{
			Name: "actor",
		},
	}
	require.Equal(t, wantPolicy, resp.Record.Policy)
	require.Equal(t, types.PolicyMarshalingType_SHORT_YAML, resp.Record.MarshalType)
	require.NoError(t, err)
}

func TestListPolicy_NoPolicies(t *testing.T) {
	ctx := test.NewTestCtx(t)

	req := types.ListPoliciesRequest{}

	resp, err := ctx.Engine.ListPolicies(ctx, &req)

	require.NoError(t, err)
	want := &types.ListPoliciesResponse{
		Records: []*types.PolicyRecord{},
	}
	require.NoError(t, err)
	require.Equal(t, want, resp)
}
