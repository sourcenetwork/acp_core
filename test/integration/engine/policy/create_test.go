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
		Id:                "ba5162bd61996b6fb6e66ef85449f0de2e89584743df7f71577674cfb531eb25",
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
						Name:       "own",
						Expression: "owner",
						Doc:        "own doc",
					},
					{
						Name:       "read",
						Expression: "(owner + reader)",
					},
				},
				ManagementRules: []*types.ManagementRule{
					{
						Relation:   "admin",
						Expression: "owner",
					},
					{
						Relation:   "owner",
						Expression: "owner",
					},
					{
						Relation:   "reader",
						Expression: "(admin + owner)",
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
}

func TestCreatePolicy_ResourcesWithoutOwnerRelation_IsAutomaticallyAdded(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	pol := `
description: ok
name: policy
resources:
- name: file
  relations:
  - name: reader
- name: foo
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

	pol := `
description: ok
name: policy
resources:
- name: file
  relations:
  - manages:
    - deleter
    name: admin
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

	want1 := "4107b53494261acabf0109bc7e7599d63459cd91db98fedc397651d413467871"
	want2 := "2eadb005094bf4b20435b03c6fb8cace4c070eae8d88d478402663d92df92c93"
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
