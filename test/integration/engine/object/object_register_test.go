package object

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func registerObjectTestSetup(t *testing.T) *test.TestCtx {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("creator")

	policy := `
    name: policy
    resources:
      resource:
        relations:
          owner:
            types:
              - actor
    `
	action := test.CreatePolicyAction{
		Policy: policy,
	}
	action.Run(ctx)

	return ctx
}

func TestRegisterObject_RegisteringNewObjectIsSucessful(t *testing.T) {
	ctx := registerObjectTestSetup(t)

	ctx.SetPrincipal("bob")
	bob := ctx.Actors.DID("bob")

	req := types.RegisterObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("resource", "foo"),
		Metadata: metadata,
	}
	resp, err := ctx.Engine.RegisterObject(ctx, &req)

	want := &types.RegisterObjectResponse{
		Record: &types.RelationshipRecord{
			PolicyId:     ctx.State.PolicyId,
			Relationship: types.NewActorRelationship("resource", "foo", "owner", bob),
			Archived:     false,
			Metadata: &types.RecordMetadata{
				Creator: &types.Principal{
					Identifier: bob,
					Kind:       types.PrincipalKind_DID,
				},
				CreationTs: ctx.Time,
				Supplied:   metadata,
			},
		},
	}
	require.NoError(t, err)
	require.Equal(t, want, resp)

	/*
		event := &types.EventObjectRegistered{
			Actor:          did,
			PolicyId:       pol.Id,
			ObjectId:       "foo",
			ObjectResource: "resource",
		}
		testutil.AssertEventEmmited(t, ctx, event)
	*/
}

func TestRegisterObject_RegisteringObjectRegisteredToAnotherUser_ErrorsForbidden(t *testing.T) {
	ctx := registerObjectTestSetup(t)

	// Given alice as the owner of foo
	ctx.SetPrincipal("alice")
	req := types.RegisterObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("resource", "foo"),
	}
	_, err := ctx.Engine.RegisterObject(ctx, &req)
	require.NoError(t, err)

	// When bob tries to register foo
	ctx.SetPrincipal("bob")
	req = types.RegisterObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("resource", "foo"),
	}
	resp, err := ctx.Engine.RegisterObject(ctx, &req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_OPERATION_FORBIDDEN)
}

func TestRegisterObject_ReregisteringObjectOwnedByUser_ReturnsOperationForbidden(t *testing.T) {
	ctx := registerObjectTestSetup(t)

	// Given alice as the owner of foo
	ctx.SetPrincipal("alice")
	req := types.RegisterObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("resource", "foo"),
	}
	_, err := ctx.Engine.RegisterObject(ctx, &req)
	require.NoError(t, err)

	// When alice tries to register foo
	ctx.SetPrincipal("alice")
	req = types.RegisterObjectRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("resource", "foo"),
	}
	resp, err := ctx.Engine.RegisterObject(ctx, &req)
	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_OPERATION_FORBIDDEN)
}

func TestRegisterObject_RegisteringAnotherUsersArchivedObject_ReturnsOperationForbidden(t *testing.T) {
	ctx := registerObjectTestSetup(t)

	// Given alice as the previous owner of foo ctx.SetPrincipal("alice")
	_, err := ctx.Engine.RegisterObject(
		ctx,
		&types.RegisterObjectRequest{
			PolicyId: ctx.State.PolicyId,
			Object:   types.NewObject("resource", "foo"),
		},
	)
	require.NoError(t, err)
	_, err = ctx.Engine.ArchiveObject(
		ctx,
		&types.ArchiveObjectRequest{
			PolicyId: ctx.State.PolicyId,
			Object:   types.NewObject("resource", "foo"),
		},
	)
	require.NoError(t, err)

	// When bob tries to register foo
	ctx.SetPrincipal("bob")
	resp, err := ctx.Engine.RegisterObject(ctx,
		&types.RegisterObjectRequest{
			PolicyId: ctx.State.PolicyId,
			Object:   types.NewObject("resource", "foo"),
		},
	)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_OPERATION_FORBIDDEN)
}

func TestRegisterObject_RegisteringArchivedUserObject_ReturnsOperationForbidden(t *testing.T) {
	ctx := registerObjectTestSetup(t)

	// Given alice as the previous owner of foo
	ctx.SetPrincipal("alice")
	_, err := ctx.Engine.RegisterObject(
		ctx,
		&types.RegisterObjectRequest{
			PolicyId: ctx.State.PolicyId,
			Object:   types.NewObject("resource", "foo"),
		},
	)
	require.NoError(t, err)
	_, err = ctx.Engine.ArchiveObject(
		ctx,
		&types.ArchiveObjectRequest{
			PolicyId: ctx.State.PolicyId,
			Object:   types.NewObject("resource", "foo"),
		},
	)
	require.NoError(t, err)

	// When alice attempt to reregister Foo
	ctx.SetPrincipal("alice")
	got, err := ctx.Engine.RegisterObject(
		ctx,
		&types.RegisterObjectRequest{
			PolicyId: ctx.State.PolicyId,
			Object:   types.NewObject("resource", "foo"),
		},
	)
	require.Nil(t, got)
	require.ErrorIs(t, err, errors.ErrorType_OPERATION_FORBIDDEN)
	/*
		event := &types.EventObjectRegistered{
			Actor:          bobDID,
			PolicyId:       pol.Id,
			ObjectId:       "foo",
			ObjectResource: "resource",
		}
		testutil.AssertEventEmmited(t, ctx, event)
	*/
}

func TestRegisterObject_RegisteringObjectInAnUndefinedResourceErrors(t *testing.T) {
	ctx := registerObjectTestSetup(t)

	ctx.SetPrincipal("alice")
	resp, err := ctx.Engine.RegisterObject(
		ctx,
		&types.RegisterObjectRequest{
			PolicyId: ctx.State.PolicyId,
			Object:   types.NewObject("undefined-resource", "foo"),
		},
	)

	require.Nil(t, resp)
	require.NotNil(t, err) // Error should be issue by Zanzi, the internal error codes aren't stable yet
}

func TestRegisterObject_RegisteringToUnknownPolicyReturnsError(t *testing.T) {
	ctx := registerObjectTestSetup(t)

	ctx.SetPrincipal("alice")
	resp, err := ctx.Engine.RegisterObject(
		ctx,
		&types.RegisterObjectRequest{
			PolicyId: "abc1234",
			Object:   types.NewObject("resource", "foo"),
		},
	)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_NOT_FOUND)
}

func TestRegisterObject_BlankResourceErrors(t *testing.T) {
	ctx := registerObjectTestSetup(t)

	ctx.SetPrincipal("alice")
	resp, err := ctx.Engine.RegisterObject(
		ctx,
		&types.RegisterObjectRequest{
			PolicyId: ctx.State.PolicyId,
			Object:   types.NewObject("", "foo"),
		},
	)

	require.Nil(t, resp)
	require.NotNil(t, err)
}

func TestRegisterObject_BlankObjectIdErrors(t *testing.T) {
	ctx := registerObjectTestSetup(t)

	ctx.SetPrincipal("alice")
	resp, err := ctx.Engine.RegisterObject(
		ctx,
		&types.RegisterObjectRequest{
			PolicyId: ctx.State.PolicyId,
			Object:   types.NewObject("resource", ""),
		},
	)

	require.Nil(t, resp)
	require.NotNil(t, err)
}
