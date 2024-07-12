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
		PolicyId:     ctx.State.PolicyId,
		Object:       types.NewObject("resource", "foo"),
		CreationTime: timestamp,
		Metadata:     metadata,
	}
	resp, err := ctx.Engine.RegisterObject(ctx, &req)

	want := &types.RegisterObjectResponse{
		Result: types.RegistrationResult_Registered,
		Record: &types.RelationshipRecord{
			PolicyId:     ctx.State.PolicyId,
			OwnerDid:     bob,
			Relationship: types.NewActorRelationship("resource", "foo", "owner", bob),
			Archived:     false,
			CreationTime: timestamp,
			Metadata:     metadata,
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

func TestRegisterObject_RegisteringObjectRegisteredToAnotherUserErrors(t *testing.T) {
	ctx := registerObjectTestSetup(t)

	// Given alice as the owner of foo
	ctx.SetPrincipal("alice")
	req := types.RegisterObjectRequest{
		PolicyId:     ctx.State.PolicyId,
		Object:       types.NewObject("resource", "foo"),
		CreationTime: timestamp,
	}
	_, err := ctx.Engine.RegisterObject(ctx, &req)
	require.NoError(t, err)

	// When bob tries to register foo
	ctx.SetPrincipal("bob")
	req = types.RegisterObjectRequest{
		PolicyId:     ctx.State.PolicyId,
		Object:       types.NewObject("resource", "foo"),
		CreationTime: timestamp,
	}
	resp, err := ctx.Engine.RegisterObject(ctx, &req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_UNAUTHORIZED)
}

func TestRegisterObject_ReregisteringObjectOwnedByUserIsNoop(t *testing.T) {
	ctx := registerObjectTestSetup(t)

	// Given alice as the owner of foo
	ctx.SetPrincipal("alice")
	req := types.RegisterObjectRequest{
		PolicyId:     ctx.State.PolicyId,
		Object:       types.NewObject("resource", "foo"),
		CreationTime: timestamp,
	}
	_, err := ctx.Engine.RegisterObject(ctx, &req)
	require.NoError(t, err)

	// When alice tries to register foo
	ctx.SetPrincipal("alice")
	alice := ctx.Actors.DID("alice")
	req = types.RegisterObjectRequest{
		PolicyId:     ctx.State.PolicyId,
		Object:       types.NewObject("resource", "foo"),
		CreationTime: timestamp,
	}
	resp, err := ctx.Engine.RegisterObject(ctx, &req)

	want := &types.RegisterObjectResponse{
		Result: types.RegistrationResult_NoOp,
		Record: &types.RelationshipRecord{
			CreationTime: timestamp,
			PolicyId:     ctx.State.PolicyId,
			Relationship: types.NewActorRelationship("resource", "foo", "owner", alice),
			Archived:     false,
			OwnerDid:     alice,
		},
	}
	require.Equal(t, want, resp)
	require.NoError(t, err)
}

func TestRegisterObject_RegisteringAnotherUsersArchivedObjectErrors(t *testing.T) {
	ctx := registerObjectTestSetup(t)

	// Given alice as the previous owner of foo
	ctx.SetPrincipal("alice")
	_, err := ctx.Engine.RegisterObject(
		ctx,
		&types.RegisterObjectRequest{
			PolicyId:     ctx.State.PolicyId,
			Object:       types.NewObject("resource", "foo"),
			CreationTime: timestamp,
		},
	)
	require.NoError(t, err)
	_, err = ctx.Engine.UnregisterObject(
		ctx,
		&types.UnregisterObjectRequest{
			PolicyId: ctx.State.PolicyId,
			Object:   types.NewObject("resource", "foo"),
		},
	)
	require.NoError(t, err)

	// When bob tries to register foo
	ctx.SetPrincipal("bob")
	resp, err := ctx.Engine.RegisterObject(ctx,
		&types.RegisterObjectRequest{
			PolicyId:     ctx.State.PolicyId,
			Object:       types.NewObject("resource", "foo"),
			CreationTime: timestamp,
		},
	)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_UNAUTHORIZED)
}

func TestRegisterObject_RegisteringArchivedUserObjectUnarchivesObject(t *testing.T) {
	ctx := registerObjectTestSetup(t)

	// Given alice as the previous owner of foo
	ctx.SetPrincipal("alice")
	_, err := ctx.Engine.RegisterObject(
		ctx,
		&types.RegisterObjectRequest{
			PolicyId:     ctx.State.PolicyId,
			Object:       types.NewObject("resource", "foo"),
			CreationTime: timestamp,
		},
	)
	require.NoError(t, err)
	_, err = ctx.Engine.UnregisterObject(
		ctx,
		&types.UnregisterObjectRequest{
			PolicyId: ctx.State.PolicyId,
			Object:   types.NewObject("resource", "foo"),
		},
	)
	require.NoError(t, err)

	// When alice attempt to reregister Foo
	ctx.SetPrincipal("alice")
	alice := ctx.Actors.DID("alice")
	got, err := ctx.Engine.RegisterObject(
		ctx,
		&types.RegisterObjectRequest{
			PolicyId:     ctx.State.PolicyId,
			Object:       types.NewObject("resource", "foo"),
			CreationTime: timestamp,
		},
	)

	want := &types.RegisterObjectResponse{
		Result: types.RegistrationResult_Unarchived,
		Record: &types.RelationshipRecord{
			CreationTime: timestamp,
			PolicyId:     ctx.State.PolicyId,
			Relationship: types.NewActorRelationship("resource", "foo", "owner", alice),
			Archived:     false,
			OwnerDid:     alice,
		},
	}
	require.Equal(t, want, got)
	require.Nil(t, err)

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
			PolicyId:     ctx.State.PolicyId,
			Object:       types.NewObject("undefined-resource", "foo"),
			CreationTime: timestamp,
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
			PolicyId:     "abc1234",
			Object:       types.NewObject("resource", "foo"),
			CreationTime: timestamp,
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
			PolicyId:     ctx.State.PolicyId,
			Object:       types.NewObject("", "foo"),
			CreationTime: timestamp,
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
			PolicyId:     ctx.State.PolicyId,
			Object:       types.NewObject("resource", ""),
			CreationTime: timestamp,
		},
	)

	require.Nil(t, resp)
	require.NotNil(t, err)
}
