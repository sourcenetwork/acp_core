package policy

import (
	"testing"

	"github.com/sourcenetwork/acp_core/internal/specification"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
	"github.com/stretchr/testify/require"
)

func TestCreatePolicy_DefraSpec_RequiresRead(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	pol := `
	name: test
	spec: defra
	resources:
	  file:
	    permissions:
		  write:
		    expr: owner
	`

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, specification.ErrDefraSpec)
}

func TestCreatePolicy_DefraSpec_RequiresWrite(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	pol := `
	name: test
	spec: defra
	resources:
	  file:
	    permissions:
		  read:
		    expr: owner
	`

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, specification.ErrDefraSpec)
}

func TestCreatePolicy_DefraSpec_OkWithReadAndWrite(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	pol := `
	name: test
	spec: defra
	resources:
	  file:
	    permissions:
		  write:
		    expr: owner
		  read:
		    expr: owner
	`

	req := types.CreatePolicyRequest{
		Policy:      pol,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	}
	resp, err := ctx.Engine.CreatePolicy(ctx, &req)

	require.NoError(t, err)
	require.NotNil(t, resp.Record.Policy)
}

func TestCreatePolicy_DefraSpec_WriteImpliesRead(t *testing.T) {
	ctx := test.NewTestCtx(t)

	pol := `
	name: test
	spec: defra
	resources:
	  file:
	    relations:
		  writer:
		    types:
			  - actor
	    permissions:
		  write:
		    expr: owner + writer
		  read:
		    expr: owner 
	`

	action := test.PolicySetupAction{
		Policy:        pol,
		PolicyCreator: "bob",
		ObjectsPerActor: map[string][]*types.Object{
			"alice": []*types.Object{
				types.NewObject("file", "readme"),
			},
		},
		RelationshipsPerActor: map[string][]*types.Relationship{
			"alice": {
				types.NewActorRelationship("file", "readme", "writer", ctx.Actors.DID("bob")),
			},
		},
	}
	action.Run(ctx)

	resp, err := ctx.Engine.VerifyAccessRequest(ctx, &types.VerifyAccessRequestRequest{
		PolicyId: ctx.State.PolicyId,
		AccessRequest: &types.AccessRequest{
			Operations: []*types.Operation{
				&types.Operation{
					Object:     types.NewObject("file", "readme"),
					Permission: "read",
				},
			},
			Actor: types.NewActor(ctx.Actors.DID("bob")),
		},
	})
	require.NoError(t, err)
	require.True(t, resp.Valid)
}
