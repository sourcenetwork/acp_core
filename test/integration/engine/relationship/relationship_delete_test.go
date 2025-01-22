package relationship

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func testDeleteRelationshipSetup(t *testing.T) *test.TestCtx {
	ctx := test.NewTestCtx(t)
	admin := ctx.Actors.DID("admin")
	writer := ctx.Actors.DID("writer")
	reader := ctx.Actors.DID("reader")

	policy := `
    name: policy
    resources:
      file:
        relations:
          owner:
            types:
              - actor
          reader:
            types:
              - actor
          writer:
            types:
              - actor
          admin:
            types:
              - actor
            manages:
              - reader
    `
	action := test.PolicySetupAction{
		Policy:        policy,
		PolicyCreator: "creator",
		ObjectsPerActor: map[string][]*types.Object{
			"alice": {
				types.NewObject("file", "foo"),
			},
		},
		RelationshipsPerActor: map[string][]*types.Relationship{
			"alice": {
				types.NewActorRelationship("file", "foo", "reader", reader),
				types.NewActorRelationship("file", "foo", "writer", writer),
				types.NewActorRelationship("file", "foo", "admin", admin),
			},
		},
	}
	action.Run(ctx)

	return ctx
}

func TestDeleteRelationship_ObjectOwnerCanRemoveRelationship(t *testing.T) {
	ctx := testDeleteRelationshipSetup(t)

	ctx.SetPrincipal("alice")
	reader := ctx.Actors.DID("reader")
	req := &types.DeleteRelationshipRequest{
		PolicyId:     ctx.State.PolicyId,
		Relationship: types.NewActorRelationship("file", "foo", "reader", reader),
	}
	resp, err := ctx.Engine.DeleteRelationship(ctx, req)

	want := &types.DeleteRelationshipResponse{
		RecordFound: true,
	}
	require.Equal(t, want, resp)
	require.NoError(t, err)
}

func TestDeleteRelationship_ObjectManagerCanRemoveRelationshipsForRelationTheyManage(t *testing.T) {
	ctx := testDeleteRelationshipSetup(t)

	ctx.SetPrincipal("admin")
	reader := ctx.Actors.DID("reader")
	req := &types.DeleteRelationshipRequest{
		PolicyId:     ctx.State.PolicyId,
		Relationship: types.NewActorRelationship("file", "foo", "reader", reader),
	}
	resp, err := ctx.Engine.DeleteRelationship(ctx, req)

	want := &types.DeleteRelationshipResponse{
		RecordFound: true,
	}
	require.Equal(t, want, resp)
	require.NoError(t, err)
}
func TestDeleteRelationship_ObjectManagerCannotRemoveRelationshipForRelationTheyDontManage(t *testing.T) {
	ctx := testDeleteRelationshipSetup(t)

	ctx.SetPrincipal("admin")
	writer := ctx.Actors.DID("writer")
	req := &types.DeleteRelationshipRequest{
		PolicyId:     ctx.State.PolicyId,
		Relationship: types.NewActorRelationship("file", "foo", "writer", writer),
	}
	resp, err := ctx.Engine.DeleteRelationship(ctx, req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_UNAUTHORIZED)
}
