package policy

import (
	"testing"

	"github.com/sourcenetwork/acp_core/internal/policy/ppp"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
	"github.com/stretchr/testify/require"
)

func TestValidatePolicy_ValidPolicyOk(t *testing.T) {
	ctx := test.NewTestCtx(t)

	pol := `
name: test
resources:
- name: foo
  permissions:
  - expr: reader
    name: read
  relations:
  - name: reader
`
	resp, err := ctx.Engine.ValidatePolicy(ctx, &types.ValidatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_YAML,
	})
	require.NoError(t, err)
	require.True(t, resp.Valid)
}

func TestValidatePolicy_InvalidPolicyReturnsErrorMsg(t *testing.T) {
	ctx := test.NewTestCtx(t)

	pol := `
name: test
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

	policyStr := `
actor:
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
    name: own
  - expr: reader
    name: read
  relations:
  - manages:
    - reader
    name: admin
  - name: reader
`

	msg := types.ValidatePolicyRequest{
		Policy:      policyStr,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp, err := ctx.Engine.ValidatePolicy(ctx, &msg)

	require.Nil(t, err)
	require.Empty(t, resp.ErrorMsg)
	require.Equal(t, &types.Policy{
		Id:                "",
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
						Name:    "reader",
						VrTypes: []*types.Restriction{},
					},
				},
				Permissions: []*types.Permission{
					{
						Name:                "own",
						Expression:          "",
						EffectiveExpression: "owner",
						Doc:                 "own doc",
					},
					{
						Name:                "read",
						Expression:          "reader",
						EffectiveExpression: "(owner + reader)",
					},
				},
				ManagementRules: []*types.ManagementRule{
					{
						Relation:   "admin",
						Expression: "owner",
						Managers:   []string{"owner"},
					},
					{
						Relation:   "owner",
						Expression: "owner",
						Managers:   []string{"owner"},
					},
					{
						Relation:   "reader",
						Expression: "(admin + owner)",
						Managers: []string{
							"admin", "owner",
						},
					},
				},
				Owner: &types.Relation{
					Name: "owner",
					Doc:  ppp.OwnerDescription,
					VrTypes: []*types.Restriction{
						{
							ResourceName: "actor-resource",
						},
					},
					Manages: []string{
						"admin", "reader", "owner",
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
		resp.Policy,
	)
}
