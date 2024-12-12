package relationship

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func setRelationshipTestSetup(t *testing.T) *test.TestCtx {
	ctx := test.NewTestCtx(t)

	policy := `
    name: policy
    resources:
      file:
        relations:
          owner:
            types:
              - actor
          admin:
            manages:
              - reader
            types:
              - actor
          reader:
            types:
              - actor
    `
	action := test.PolicySetupAction{
		Policy:        policy,
		PolicyCreator: "root",
		ObjectsPerActor: map[string][]*types.Object{
			"alice": []*types.Object{
				types.NewObject("file", "foo"),
			},
		},
	}
	action.Run(ctx)

	return ctx
}

func TestSetRelationship_OwnerCanShareObjectTheyOwn(t *testing.T) {
	ctx := setRelationshipTestSetup(t)

	ctx.SetPrincipal("alice")
	alice := ctx.Actors.DID("alice")
	bob := ctx.Actors.DID("bob")
	req := &types.SetRelationshipRequest{
		PolicyId:     ctx.State.PolicyId,
		Relationship: types.NewActorRelationship("file", "foo", "reader", bob),
		Metadata:     metadata,
	}
	resp, err := ctx.Engine.SetRelationship(ctx, req)

	want := &types.SetRelationshipResponse{
		RecordExisted: false,
		Record: &types.RelationshipRecord{
			PolicyId:     ctx.State.PolicyId,
			Relationship: types.NewActorRelationship("file", "foo", "reader", bob),
			Archived:     false,
			Metadata: &types.RecordMetadata{
				CreationTs: ctx.Time,
				Creator: &types.Principal{
					Kind:       types.PrincipalKind_DID,
					Identifier: alice,
				},
				Supplied: metadata,
			},
		},
	}
	require.Equal(t, want, resp)
	require.NoError(t, err)
}
func TestSetRelationship_ActorCannotSetRelationshipForUnregisteredObject(t *testing.T) {
	ctx := setRelationshipTestSetup(t)

	bob := ctx.Actors.DID("bob")
	ctx.SetPrincipal("alice")
	req := &types.SetRelationshipRequest{
		PolicyId:     ctx.State.PolicyId,
		Relationship: types.NewActorRelationship("file", "unregistered", "reader", bob),
	}
	resp, err := ctx.Engine.SetRelationship(ctx, req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_NOT_FOUND)
}

func TestSetRelationship_ActorCannotSetRelationshipForObjectTheyDoNotOwn(t *testing.T) {
	ctx := setRelationshipTestSetup(t)

	ctx.SetPrincipal("bob")
	bob := ctx.Actors.DID("bob")
	req := &types.SetRelationshipRequest{
		PolicyId:     ctx.State.PolicyId,
		Relationship: types.NewActorRelationship("file", "foo", "reader", bob),
	}
	resp, err := ctx.Engine.SetRelationship(ctx, req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_UNAUTHORIZED)
}

func TestSetRelationship_ManagerActorCanDelegateAccessToAnotherActor(t *testing.T) {
	ctx := setRelationshipTestSetup(t)

	// Given object foo and Bob as a manager
	ctx.SetPrincipal("alice")
	bob := ctx.Actors.DID("bob")
	a1 := test.SetRelationshipsAction{
		PolicyId: ctx.State.PolicyId,
		Relationships: []*types.Relationship{
			types.NewActorRelationship("file", "foo", "admin", bob),
		},
	}
	a1.Run(ctx)

	// when bob shares foo with charlie
	ctx.SetPrincipal("bob")
	charlie := ctx.Actors.DID("charlie")
	req := &types.SetRelationshipRequest{
		PolicyId:     ctx.State.PolicyId,
		Relationship: types.NewActorRelationship("file", "foo", "reader", charlie),
	}
	resp, err := ctx.Engine.SetRelationship(ctx, req)

	want := &types.SetRelationshipResponse{
		RecordExisted: false,
		Record: &types.RelationshipRecord{
			PolicyId:     ctx.State.PolicyId,
			Relationship: types.NewActorRelationship("file", "foo", "reader", charlie),
			Archived:     false,
			Metadata: &types.RecordMetadata{
				Creator: &types.Principal{
					Identifier: bob,
					Kind:       types.PrincipalKind_DID,
				},
				CreationTs: ctx.Time,
				Supplied:   nil,
			},
		},
	}
	require.Equal(t, want, resp)
	require.NoError(t, err)
}

func TestSetRelationship_ManagerActorCannotSetRelationshipToRelationshipsTheyDoNotManage(t *testing.T) {
	ctx := setRelationshipTestSetup(t)

	// Given object foo and Bob as a admin
	ctx.SetPrincipal("alice")
	bob := ctx.Actors.DID("bob")
	a1 := test.SetRelationshipsAction{
		PolicyId: ctx.State.PolicyId,
		Relationships: []*types.Relationship{
			types.NewActorRelationship("file", "foo", "admin", bob),
		},
	}
	a1.Run(ctx)

	// when bob attemps to make charlie an admin
	ctx.SetPrincipal("bob")
	charlie := ctx.Actors.DID("charlie")
	req := &types.SetRelationshipRequest{
		PolicyId:     ctx.State.PolicyId,
		Relationship: types.NewActorRelationship("file", "foo", "admin", charlie),
	}
	resp, err := ctx.Engine.SetRelationship(ctx, req)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_UNAUTHORIZED)
}
