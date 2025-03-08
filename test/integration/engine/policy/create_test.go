package policy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/internal/policy/ppp"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

var metadata *types.SuppliedMetadata = &types.SuppliedMetadata{
	Attributes: map[string]string{
		"test": "abc",
	},
	Blob: []byte("test"),
}

func TestCreatePolicy_ValidPolicyIsCreated(t *testing.T) {
	ctx := test.NewTestCtx(t)

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

meta:
  a: b
  key: value

actor:
  name: actor-resource
  doc: my actor
`
	msg := types.CreatePolicyRequest{
		Policy:      policyStr,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
		Metadata:    metadata,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &msg)

	require.Nil(t, err)
	p := types.AnonymousPrincipal()
	wantMetadata := &types.RecordMetadata{
		Creator:    &p,
		CreationTs: ctx.Time,
		Supplied:   metadata,
	}
	require.Equal(t, wantMetadata, resp.Record.Metadata)
	require.Equal(t, &types.Policy{
		Id:          "7c10ce44694f17560cf9f5b310f90ca3db722007bae4d82e484b7e1a9fc1294b",
		Name:        "policy",
		Description: "ok",
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
						Name: "reader",
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
						Expression: "(owner + admin)",
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
			Name: "actor-resource",
			Doc:  "my actor",
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

	require.NoError(t, err)
	want := &types.Relation{
		Name:    "owner",
		Doc:     "owner relations represents the object owner",
		Manages: nil,
		VrTypes: []*types.Restriction{
			&types.Restriction{
				ResourceName: resp.Record.Policy.ActorResource.Name,
			},
		},
	}
	require.Equal(t, want, resp.Record.Policy.Resources[0].Relations[0])
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
	require.ErrorIs(t, err, ppp.ErrAdministrationTransformer)
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
	require.ErrorIs(t, err, errors.ErrInvalidPolicy)
}

func TestCreatePolicy_CreatingMultipleEqualPoliciesProduceDifferentIDs(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("creator")

	pol := `
name: test
description: A Valid Defra Policy Interface (DPI)

actor:
  name: actor

resources:
  users:
    permissions:
      read:
        expr: owner + reader
      write:
        expr: owner

    relations:
      owner:
        types:
          - actor
      reader:
        types:
          - actor
      admin:
        manages:
          - reader
        types:
          - actor
`

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	}
	resp1, err1 := ctx.Engine.CreatePolicy(ctx, &req)
	resp2, err2 := ctx.Engine.CreatePolicy(ctx, &req)

	want1 := "60c328483733cf9607b7eca3c270784ebe25038f2c1cc0f3ae9beddeb6b1acb3"
	want2 := "4e3ce429d500c85fb9b5439ab31d2276f014a689e3ee7055cecd6e510654f4f9"
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.Equal(t, want1, resp1.Record.Policy.Id)
	require.Equal(t, want2, resp2.Record.Policy.Id)
}
