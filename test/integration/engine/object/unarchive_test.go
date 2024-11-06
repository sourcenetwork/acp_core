package object

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func unarchiveSetup(t *testing.T) *test.TestCtx {
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
	action := test.PolicySetupAction{
		Policy:        policy,
		PolicyCreator: "creator",
		ObjectsPerActor: map[string][]*types.Object{
			"bob": []*types.Object{
				types.NewObject("resource", "active"),
				types.NewObject("resource", "archived"),
			},
		},
	}
	action.Run(ctx)

	a := test.ArchiveObjectAction{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("resource", "archived"),
	}
	a.Run(ctx)

	return ctx
}

func TestUnarchive_UnarchiveUnregisteredObject_Errors(t *testing.T) {
	ctx := unarchiveSetup(t)

	ctx.SetPrincipal("bob")
	a := test.UnarchiveObjectAction{
		PolicyId:    ctx.State.PolicyId,
		Object:      types.NewObject("resource", "not-registered"),
		ExpectedErr: errors.ErrorType_BAD_INPUT,
	}
	a.Run(ctx)
}

func TestUnarchive_ActiveObjectOwnedByActor_NoopAndModifiedFalse(t *testing.T) {
	ctx := unarchiveSetup(t)

	ctx.SetPrincipal("bob")
	bob := ctx.Actors.DID("bob")
	a := test.UnarchiveObjectAction{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("resource", "active"),
		Expected: &types.UnarchiveObjectResponse{
			RecordModified: false,
			Record: &types.RelationshipRecord{
				PolicyId:     ctx.State.PolicyId,
				OwnerDid:     bob,
				Archived:     false,
				CreationTime: ctx.Time,
				Relationship: types.NewActorRelationship("resource", "active", "owner", bob),
			},
		},
	}
	a.Run(ctx)
}

func TestUnarchive_ActiveObjectOwnedBySomeoneElse_Unauthorized(t *testing.T) {
	ctx := unarchiveSetup(t)

	ctx.SetPrincipal("alice")
	a := test.UnarchiveObjectAction{
		PolicyId:    ctx.State.PolicyId,
		Object:      types.NewObject("resource", "active"),
		ExpectedErr: errors.ErrorType_UNAUTHORIZED,
	}
	a.Run(ctx)
}

func TestUnarchive_UnarchiveActiveObjectOwnedBySomeoneElse_Errors(t *testing.T) {
	ctx := unarchiveSetup(t)

	ctx.SetPrincipal("alice")
	a := test.UnarchiveObjectAction{
		PolicyId:    ctx.State.PolicyId,
		Object:      types.NewObject("resource", "active"),
		ExpectedErr: errors.ErrorType_UNAUTHORIZED,
	}
	a.Run(ctx)
}

func TestUnarchive_UnarchiveMyOwnObject_Ok(t *testing.T) {
	ctx := unarchiveSetup(t)

	ctx.SetPrincipal("bob")
	bob := ctx.Actors.DID("bob")

	a := test.UnarchiveObjectAction{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("resource", "archived"),
		Expected: &types.UnarchiveObjectResponse{
			RecordModified: true,
			Record: &types.RelationshipRecord{
				PolicyId:     ctx.State.PolicyId,
				OwnerDid:     bob,
				Archived:     false,
				CreationTime: ctx.Time,
				Relationship: types.NewActorRelationship("resource", "archived", "owner", bob),
			},
		},
	}
	a.Run(ctx)
}
func TestUnarchive_SomeoneElsesArchivedObject_Errors(t *testing.T) {
	ctx := unarchiveSetup(t)

	ctx.SetPrincipal("alice")
	a := test.UnarchiveObjectAction{
		PolicyId:    ctx.State.PolicyId,
		Object:      types.NewObject("resource", "archived"),
		ExpectedErr: errors.ErrorType_UNAUTHORIZED,
	}
	a.Run(ctx)
}
