package policy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func TestCreatePolicyWithSpec_ValidPolicyIsCreated(t *testing.T) {
	ctx := test.NewTestCtx(t)

	policyStr := `
name: policy
description: ok
spec: none
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
	msg := types.CreatePolicyWithSpecificationRequest{
		Policy:       policyStr,
		MarshalType:  types.PolicyMarshalingType_SHORT_YAML,
		Metadata:     metadata,
		RequiredSpec: types.PolicySpecificationType_NO_SPEC,
	}
	resp, err := ctx.Engine.CreatePolicyWithSpecification(ctx, &msg)

	require.Nil(t, err)
	p := types.AnonymousPrincipal()
	wantMetadata := &types.RecordMetadata{
		Creator:    &p,
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
			Name: "actor-resource",
			Doc:  "my actor",
		},
	},
		resp.Record.Policy,
	)
}

func TestCreatePolicyWithSpec_RequiredSpecDiffersFromInformedSpec(t *testing.T) {
	ctx := test.NewTestCtx(t)

	policyStr := `
name: policy
spec: none
`
	msg := types.CreatePolicyWithSpecificationRequest{
		Policy:       policyStr,
		MarshalType:  types.PolicyMarshalingType_SHORT_YAML,
		RequiredSpec: types.PolicySpecificationType_DEFRA_SPEC,
	}
	resp, err := ctx.Engine.CreatePolicyWithSpecification(ctx, &msg)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_BAD_INPUT)
}
