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
		Id:                "199091661bdd06221eb0a8070673c76f25ca8c8dcc04d47934f0abb123daf78b",
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
				Owner: &types.Relation{
					Name: "owner",
					Doc:  ppp.OwnerDescription,
					VrTypes: []*types.Restriction{
						{
							ResourceName: "actor",
						},
					},
					Manages: []string{
						"admin",
						"reader",
						"owner",
					},
				},
				ManagementRules: []*types.ManagementRule{
					{
						Relation:   "admin",
						Expression: "owner",
						Managers: []string{
							"owner",
						},
					},
					{
						Relation:   "owner",
						Expression: "owner",
						Managers: []string{
							"owner",
						},
					},
					{
						Relation:   "reader",
						Expression: "(admin + owner)",
						Managers: []string{
							"admin",
							"owner",
						},
					},
				},
			},
		},
		ActorResource: &types.ActorResource{
			Name: ppp.ActorResourceName,
			Doc:  ppp.ActorResourceDoc,
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
		Name: "owner",
		Doc:  ppp.OwnerDescription,
		Manages: []string{
			"reader",
			"owner",
		},
		VrTypes: []*types.Restriction{
			{
				ResourceName: resp.Record.Policy.ActorResource.Name,
			},
		},
	}
	require.Equal(t, want, resp.Record.Policy.Resources[0].Owner)
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

	pol := `
description: A Valid Defra Policy Interface (DPI)
name: test
resources:
- name: users
  permissions:
  - expr: reader
    name: read
  - name: write
  relations:
  - manages:
    - reader
    name: admin
    types:
    - actor
  - name: reader
    types:
    - actor
`

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp1, err1 := ctx.Engine.CreatePolicy(ctx, &req)
	resp2, err2 := ctx.Engine.CreatePolicy(ctx, &req)

	want1 := "9372b5ef92b9332f597b46120026583a3ceed09d50046da159ee65273602fa82"
	want2 := "4bd4dd14eec0fb91e3eb2e2f8d36b52b521bf72b60ec137658804dac8e69379e"
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

func TestCreatePolicy_ClashingRelationAndPermissionNames_Errors(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	pol := `
name: policy
resources:
- name: foo
  permissions:
  - name: test
  relations:
  - name: test
`
	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &req)
	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_BAD_INPUT)
}

func TestCreatePolicy_WithExplicitActorResource_Errors(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	pol := `
name: policy
resources:
- name: actor
  permissions:
  - name: test
`

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &req)
	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrInvalidPolicy)
}
