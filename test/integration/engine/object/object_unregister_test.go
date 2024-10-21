package object

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func testUnregisterObjectSetup(t *testing.T) *test.TestCtx {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("admin")

	pol := `
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
    `
	action := test.CreatePolicyAction{
		Policy: pol,
	}
	policy := action.Run(ctx)

	ctx.SetPrincipal("alice")
	alice := ctx.Actors.DID("alice")
	regObjAction := test.RegisterObjectsAction{
		PolicyId: policy.Id,
		Objects: []*types.Object{
			types.NewObject("file", "foo"),
		},
	}
	regObjAction.Run(ctx)

	relsAction := test.SetRelationshipsAction{
		PolicyId: policy.Id,
		Relationships: []*types.Relationship{
			types.NewActorRelationship("file", "foo", "reader", alice),
		},
	}
	relsAction.Run(ctx)

	return ctx
}

func TestUnregisterObject_RegisteredObjectCanBeUnregisteredByAuthor(t *testing.T) {
	ctx := testUnregisterObjectSetup(t)
	ctx.SetPrincipal("alice")

	req := &types.ArchiveObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("file", "foo"),
	}
	resp, err := ctx.Engine.UnregisterObject(ctx, req)

	want := &types.ArchiveObjectResponse{
		Found:                true,
		RelationshipsRemoved: 2,
	}
	require.Equal(t, want, resp)
	require.NoError(t, err)
}

func TestUnregisterObject_ActorCannotUnregisterObjectTheyDoNotOwn(t *testing.T) {
	ctx := testUnregisterObjectSetup(t)
	ctx.SetPrincipal("bob")

	req := &types.ArchiveObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("file", "foo"),
	}
	resp, err := ctx.Engine.UnregisterObject(ctx, req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_UNAUTHORIZED)
}

func TestUnregisterObject_UnregisteringAnObjectThatDoesNotExistReturnsFoundFalse(t *testing.T) {
	ctx := testUnregisterObjectSetup(t)
	ctx.SetPrincipal("alice")

	req := &types.ArchiveObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("file", "file-that-isn't-registered"),
	}
	resp, err := ctx.Engine.UnregisterObject(ctx, req)

	require.Equal(t, &types.ArchiveObjectResponse{
		Found: false,
	}, resp)
	require.NoError(t, err, errors.ErrorType_UNAUTHORIZED)
}

func TestUnregisterObject_UnregisteringAnAlreadyArchivedObjectIsANoop(t *testing.T) {
	ctx := testUnregisterObjectSetup(t)

	// Given the file Foo archived by alice
	ctx.SetPrincipal("alice")
	_, err := ctx.Engine.UnregisterObject(ctx, &types.ArchiveObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("file", "foo"),
	})
	require.NoError(t, err)

	// When alice file foo
	ctx.SetPrincipal("alice")
	resp, err := ctx.Engine.UnregisterObject(ctx, &types.ArchiveObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("file", "foo"),
	})

	want := &types.ArchiveObjectResponse{
		Found: true,
	}
	require.Equal(t, want, resp)
	require.NoError(t, err)
}

func TestUnregisterObject_SendingInvalidPolicyIdErrors(t *testing.T) {
	ctx := testUnregisterObjectSetup(t)

	// When alice file foo
	ctx.SetPrincipal("alice")
	resp, err := ctx.Engine.UnregisterObject(ctx, &types.ArchiveObjectRequest{
		PolicyId: "invalid-policy-id",
		Object:   types.NewObject("file", "foo"),
	})

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_NOT_FOUND)
}

/*
func TestUnregisterObject_UnregisteringObjectRemovesRelationshipsLeavingTheObject(t *testing.T) {
	// TODO
}
*/
