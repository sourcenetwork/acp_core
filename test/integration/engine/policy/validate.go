package policy

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
	"github.com/stretchr/testify/require"
)

func TestValidatePolicy_ValidPolicyOk(t *testing.T) {
	ctx := test.NewTestCtx(t)

	pol := `name: test
resources:
- name: foo
  permissions:
  - expr: reader
    name: read
  relations:
  - name: reader
spec: none
`
	resp, err := ctx.Engine.ValidatePolicy(ctx, &types.ValidatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_YAML,
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

	pol := `name: test
resources:
- name: foo
  permissions:
  - expr: reader
    name: read
  relations:
  - name: reader
spec: defra
`

	resp, err := ctx.Engine.ValidatePolicy(ctx, &types.ValidatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_YAML,
	})
	require.NoError(t, err)
	require.False(t, resp.Valid)
	require.Contains(t, resp.ErrorMsg, "defra policy specification")
}

func TestValidatePolicy_ReturnsParsedPolicy(t *testing.T) {
	ctx := test.NewTestCtx(t)

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

	msg := types.ValidatePolicyRequest{
		Policy:      policyStr,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp, err := ctx.Engine.ValidatePolicy(ctx, &msg)

	require.Nil(t, err)
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
		resp.Policy,
	)
}
