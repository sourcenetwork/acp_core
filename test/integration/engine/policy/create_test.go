package policy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/internal/policy/ppp"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func TestCreatePolicy_ValidPolicyIsCreated(t *testing.T) {
	ctx := test.NewTestCtx(t)
	bob := ctx.SetPrincipal("bob")

	policyStr := `actor:
  doc: my actor
  name: actor-resource
description: ok
meta:
  a: b
  key: value
name: policy
resources:
- name: file
  permissions:
  - doc: own doc
    expr: owner
    name: own
  - expr: owner + reader
    name: read
  relations:
  - manages:
    - reader
    name: admin
  - doc: owner owns
    name: owner
    types:
    - actor-resource
  - name: reader
spec: none
`

	msg := types.CreatePolicyRequest{
		Policy:      policyStr,
		MarshalType: types.PolicyMarshalingType_YAML,
		Metadata:    metadata,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &msg)

	require.Nil(t, err)
	wantMetadata := &types.RecordMetadata{
		Creator:    &bob,
		CreationTs: ctx.Time,
		Supplied:   metadata,
	}
	require.Equal(t, wantMetadata, resp.Record.Metadata)
	require.Equal(t, &types.Policy{
		Id:                "da7be65027664708551f97197ba5f5993aa99bc7b57055df9766426dc6da9605",
		Name:              "policy",
		Description:       "ok",
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Attributes: map[string]string{
			"a":   "b",
			"key": "value",
		},
		Resources: []*types.Resource{
			{
				Name: "file",
				Relations: []*types.Relation{
					{
						Name: "admin",
						Manages: []string{
							"reader",
						},
						VrTypes: []*types.Restriction{},
					},
					{
						Name: "owner",
						Doc:  "owner owns",
						VrTypes: []*types.Restriction{
							{
								ResourceName: "actor-resource",
								RelationName: "",
							},
						},
					},
					{
						Name:    "reader",
						VrTypes: []*types.Restriction{},
					},
				},
				Permissions: []*types.Permission{
					{
						Name:       "_can_manage_admin",
						Expression: "owner",
						Doc:        "permission controls actors which are allowed to create relationships for the admin relation (permission was auto-generated).",
					},
					{
						Name:       "_can_manage_owner",
						Expression: "owner",
						Doc:        "permission controls actors which are allowed to create relationships for the owner relation (permission was auto-generated).",
					},
					{
						Name:       "_can_manage_reader",
						Expression: "(admin + owner)",
						Doc:        "permission controls actors which are allowed to create relationships for the reader relation (permission was auto-generated).",
					},
					{
						Name:       "own",
						Expression: "owner",
						Doc:        "own doc",
					},
					{
						Name:       "read",
						Expression: "(owner + reader)",
					},
				},
			},
		},
		ActorResource: &types.ActorResource{
			Name:      "actor-resource",
			Doc:       "my actor",
			Relations: []*types.Relation{},
		},
	},
		resp.Record.Policy,
	)

	event := &types.EventPolicyCreated{
		PolicyId:   "4419a8abb886c641bc794b9b3289bc2118ab177542129627b6b05d540de03e46",
		PolicyName: "policy",
	}
	_ = event
	//.AssertEventEmmited(t, ctx, event)
}

func TestCreatePolicy_ResourcesWithoutOwnerRelation_IsAutomaticallyAdded(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	pol := `description: ok
name: policy
resources:
- name: file
  relations:
  - name: reader
- name: foo
  relations:
  - name: owner
spec: none
`

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &req)

	require.NoError(t, err)
	want := &types.Relation{
		Name:    "owner",
		Doc:     "owner relations represents the object owner",
		Manages: nil,
		VrTypes: []*types.Restriction{
			{
				ResourceName: resp.Record.Policy.ActorResource.Name,
			},
		},
	}
	require.Equal(t, want, resp.Record.Policy.Resources[0].Relations[0])
}

func TestCreatePolicy_ManagementReferencingUndefinedRelationReturnsError(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	pol := `description: ok
name: policy
resources:
- name: file
  relations:
  - manages:
    - deleter
    name: admin
  - name: owner
spec: none
`

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, ppp.ErrAdministrationTransformer)
}

func TestCreatePolicy_UnamedPolicyCausesError(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	pol := `spec: none
`

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrInvalidPolicy)
}

func TestCreatePolicy_CreatingMultipleEqualPoliciesProduceDifferentIDs(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("creator")

	pol := `actor:
  name: actor
description: A Valid Defra Policy Interface (DPI)
name: test
resources:
- name: users
  permissions:
  - name: read
    expr: reader
  - name: write
  relations:
  - manages:
    - reader
    name: admin
    types:
    - actor
  - name: owner
    types:
    - actor
  - name: reader
    types:
    - actor
spec: none
`

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp1, err1 := ctx.Engine.CreatePolicy(ctx, &req)
	resp2, err2 := ctx.Engine.CreatePolicy(ctx, &req)

	want1 := "0aceb40e813c157152dd931f0f5e59228fce7c87ab3a40341ac1abce7ad7da3a"
	want2 := "f3743fa4268b48462014e0b1b8a07d7a8bf615b3ddc276d91c7128aaad8a2eee"
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.Equal(t, want1, resp1.Record.Policy.Id)
	require.Equal(t, want2, resp2.Record.Policy.Id)
}

func TestCreatePolicy_WithEmptyPermission_OwnerIsPermitted(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	pol := `
name: policy
resources:
- name: foo
  relations:
  - name: owner
    types:
    - actor
  permissions:
  - name: test
`

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &req)
	require.NoError(t, err)

	_, err = ctx.Engine.RegisterObject(ctx, &types.RegisterObjectRequest{
		PolicyId: resp.Record.Policy.Id,
		Object:   types.NewObject("foo", "obj"),
	})
	require.NoError(t, err)

	checkResult, err := ctx.Engine.VerifyAccessRequest(ctx, &types.VerifyAccessRequestRequest{
		PolicyId: resp.Record.Policy.Id,
		AccessRequest: &types.AccessRequest{
			Operations: []*types.Operation{
				{
					Object:     types.NewObject("foo", "obj"),
					Permission: "test",
				},
			},
			Actor: types.NewActor(ctx.Actors.DID("bob")),
		},
	})
	require.NoError(t, err)
	require.True(t, checkResult.Valid)
}
