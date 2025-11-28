package object

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func testArchiveObjectSetup(t *testing.T) *test.TestCtx {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("admin")

	pol := `
name: policy
resources:
- name: file
  relations:
  - name: reader
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

func TestArchiveObject_RegisteredObjectCanBeArchivedByAuthor(t *testing.T) {
	ctx := testArchiveObjectSetup(t)
	ctx.SetPrincipal("alice")

	req := &types.ArchiveObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("file", "foo"),
	}
	resp, err := ctx.Engine.ArchiveObject(ctx, req)

	want := &types.ArchiveObjectResponse{
		RelationshipsRemoved: 2,
		RecordModified:       true,
	}
	require.Equal(t, want, resp)
	require.NoError(t, err)
}

func TestArchiveObject_ActorCannotArchiveObjectTheyDoNotOwn(t *testing.T) {
	ctx := testArchiveObjectSetup(t)
	ctx.SetPrincipal("bob")

	req := &types.ArchiveObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("file", "foo"),
	}
	resp, err := ctx.Engine.ArchiveObject(ctx, req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_UNAUTHORIZED)
}

func TestArchiveObject_ArchiveingAnObjectThatDoesNotExistReturnsFoundFalse(t *testing.T) {
	ctx := testArchiveObjectSetup(t)
	ctx.SetPrincipal("alice")

	req := &types.ArchiveObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("file", "file-that-isn't-registered"),
	}
	resp, err := ctx.Engine.ArchiveObject(ctx, req)
	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_BAD_INPUT)
}

func TestArchiveObject_ArchiveingAnAlreadyArchivedObjectIsANoop(t *testing.T) {
	ctx := testArchiveObjectSetup(t)

	// Given the file Foo archived by alice
	ctx.SetPrincipal("alice")
	_, err := ctx.Engine.ArchiveObject(ctx, &types.ArchiveObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("file", "foo"),
	})
	require.NoError(t, err)

	// When alice file foo
	ctx.SetPrincipal("alice")
	resp, err := ctx.Engine.ArchiveObject(ctx, &types.ArchiveObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("file", "foo"),
	})

	want := &types.ArchiveObjectResponse{
		RecordModified:       false,
		RelationshipsRemoved: 0,
	}
	require.Equal(t, want, resp)
	require.NoError(t, err)
}

func TestArchiveObject_SendingInvalidPolicyIdErrors(t *testing.T) {
	ctx := testArchiveObjectSetup(t)

	// When alice file foo
	ctx.SetPrincipal("alice")
	resp, err := ctx.Engine.ArchiveObject(ctx, &types.ArchiveObjectRequest{
		PolicyId: "invalid-policy-id",
		Object:   types.NewObject("file", "foo"),
	})

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_NOT_FOUND)
}
