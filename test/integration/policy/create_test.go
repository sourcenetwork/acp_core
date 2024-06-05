package policy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
	testutil "github.com/sourcenetwork/acp_core/test/util"
)

var timestamp = testutil.MustDateTimeToProto("2024-01-01 00:00:00")

func TestCreatePolicy_ValidPolicyIsCreated(t *testing.T) {
	ctx := test.NewTestCtx(t)
	bob := ctx.Actors.DID("bob")
	ctx.SetPrincipal("bob")

	policyStr := `
name: policy
description: ok
resources:
  file:
    relations:
      owner:
        doc: owner owns
        types:
          - actor-resource
      reader:
      admin:
        manages:
          - reader
    permissions:
      own:
        expr: owner
        doc: own doc
      read:
        expr: owner + reader
actor:
  name: actor-resource
  doc: my actor
`
	msg := types.CreatePolicyRequest{
		Policy:       policyStr,
		MarshalType:  types.PolicyMarshalingType_SHORT_YAML,
		CreationTime: timestamp,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &msg)

	require.Nil(t, err)
	require.Equal(t, resp.Policy, &types.Policy{
		Id:           "d12fa4d041911f2c77f6f49dd73942fb03389ab983714315af67b0f8e7cbcfef",
		Name:         "policy",
		Description:  "ok",
		CreationTime: timestamp,
		Creator:      bob,
		Resources: []*types.Resource{
			&types.Resource{
				Name: "file",
				Relations: []*types.Relation{
					&types.Relation{
						Name: "admin",
						Manages: []string{
							"reader",
						},
						VrTypes: []*types.Restriction{},
					},
					&types.Relation{
						Name: "owner",
						Doc:  "owner owns",
						VrTypes: []*types.Restriction{
							&types.Restriction{
								ResourceName: "actor-resource",
								RelationName: "",
							},
						},
					},
					&types.Relation{
						Name: "reader",
					},
				},
				Permissions: []*types.Permission{
					&types.Permission{
						Name:       "own",
						Expression: "owner",
						Doc:        "own doc",
					},
					&types.Permission{
						Name:       "read",
						Expression: "owner + reader",
					},
					&types.Permission{
						Name:       "_can_manage_admin",
						Expression: "owner",
						Doc:        "permission controls actors which are allowed to create relationships for the admin relation (permission was auto-generated by SourceHub).",
					},
					&types.Permission{
						Name:       "_can_manage_owner",
						Expression: "owner",
						Doc:        "permission controls actors which are allowed to create relationships for the owner relation (permission was auto-generated by SourceHub).",
					},
					&types.Permission{
						Name:       "_can_manage_reader",
						Expression: "(admin + owner)",
						Doc:        "permission controls actors which are allowed to create relationships for the reader relation (permission was auto-generated by SourceHub).",
					},
				},
			},
		},
		ActorResource: &types.ActorResource{
			Name: "actor-resource",
			Doc:  "my actor",
		},
	})

	event := &types.EventPolicyCreated{
		Creator:    bob,
		PolicyId:   "4419a8abb886c641bc794b9b3289bc2118ab177542129627b6b05d540de03e46",
		PolicyName: "policy",
	}
	_ = event
	//.AssertEventEmmited(t, ctx, event)
}

func TestCreatePolicy_PolicyResourcesRequiresOwnerRelation(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	pol := `
name: policy
description: ok
resources:
  file:
    relations:
      reader:
    permissions:
  foo:
    relations:
      owner:
    permissions:
`

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, policy.ErrResourceMissingOwnerRelation)
}

func TestCreatePolicy_ManagementReferencingUndefinedRelationReturnsError(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	pol := `
name: policy
description: ok
resources:
  file:
    relations:
      owner:
      admin:
        manages:
          - deleter
    permissions:
`

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, policy.ErrInvalidManagementRule)
}

func TestCreatePolicy_UnamedPolicyCausesError(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	pol := `
resources:
`

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, policy.ErrInvalidPolicy)
}

func TestCreatePolicy_CreatePolicyWithAnonymousPrincipalErrors(t *testing.T) {
	pol := `
name: test
resources:
`
	ctx := test.NewTestCtx(t)

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, auth.ErrUnauthenticatd)
}