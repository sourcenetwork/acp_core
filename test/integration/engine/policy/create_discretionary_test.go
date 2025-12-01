package policy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/internal/policy/ppp"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func TestCreatePolicy_ResourceContainsOwnerRelation_Errors(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	policyStr := `
name: pol
resources:
- name: file
  relations:
  - name: owner
    types:
    - actor
  permissions:
  - name: read
`

	msg := types.CreatePolicyRequest{
		Policy:      policyStr,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &msg)

	require.Nil(t, resp)
	require.ErrorIs(t, err, ppp.ErrResourceContainsOwner)
	require.ErrorIs(t, err, errors.ErrorType_BAD_INPUT)
}

func TestCreatePolicy_GeneratedOwnerRelation_CanManageEveryRelationInResource(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	policyStr := `
name: pol
resources:
- name: file
  permissions:
  - name: read
  relations:
  - name: abc
  - name: def
  - name: relation
    types:
    - actor
`

	msg := types.CreatePolicyRequest{
		Policy:      policyStr,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &msg)

	require.Nil(t, err)
	pol := resp.Record.Policy
	file := pol.GetResourceByName("file")
	want := []string{
		"abc", "def", "relation", "owner",
	}
	require.Equal(t, want, file.Owner.Manages)
}

func TestCreatePolicy_PermissionReferencingOwner_IsRejected(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	policyStr := `
name: pol
resources:
- name: file
  permissions:
  - name: read
    expr: owner
  relations:
`

	msg := types.CreatePolicyRequest{
		Policy:      policyStr,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &msg)

	require.Nil(t, resp)
	require.ErrorIs(t, err, ppp.ErrPermissionReferencesOwner)
	require.ErrorIs(t, err, errors.ErrorType_BAD_INPUT)
}

func TestCreatePolicy_PermissionReferencingAnothersResourcesOwner_IsAccepted(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	policyStr := `
name: pol
resources:
- name: directory
- name: file
  permissions:
  - name: read
    expr: parent->owner
  relations:
  - name: parent
    types:
    - directory
`

	msg := types.CreatePolicyRequest{
		Policy:      policyStr,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &msg)

	require.NoError(t, err)
	require.NotNil(t, resp)
	read := resp.Record.Policy.GetResourceByName("file").GetPermissionByName("read")
	require.Equal(t, "parent->owner", read.Expression)
}

func TestCreatePolicy_InvalidPermissionWithOwnerAsTuplesetFilter_Errors(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	policyStr := `
name: pol
resources:
- name: file
  permissions:
  - name: read
    expr: owner->something
`

	msg := types.CreatePolicyRequest{
		Policy:      policyStr,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &msg)

	require.Nil(t, resp)
	require.ErrorIs(t, err, ppp.ErrPermissionReferencesOwner)
	require.ErrorIs(t, err, errors.ErrorType_BAD_INPUT)
}

func TestCreatePolicy_WithEmptyPermissionExpression_AcceptedAndOwnerCanPerformOperation(t *testing.T) {
	ctx := test.NewTestCtx(t)

	policyStr := `
name: pol
resources:
- name: file
  permissions:
  - name: read
`

	a := test.PolicySetupAction{
		Policy:        policyStr,
		PolicyCreator: "bob",
		ObjectsPerActor: map[string][]*types.Object{
			"alice": {
				types.NewObject("file", "test"),
			},
		},
	}
	a.Run(ctx)

	// When Alice tries to read file test
	resp, err := ctx.Engine.VerifyAccessRequest(ctx, &types.VerifyAccessRequestRequest{
		PolicyId: ctx.State.PolicyId,
		AccessRequest: &types.AccessRequest{
			Operations: []*types.Operation{
				{
					Object:     types.NewObject("file", "test"),
					Permission: "read",
				},
			},
			Actor: types.NewActor(ctx.Actors.DID("alice")),
		},
	})
	require.NoError(t, err)
	require.True(t, resp.Valid)
}
