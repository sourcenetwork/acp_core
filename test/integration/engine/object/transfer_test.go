package object

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func transferSetup(t *testing.T) *test.TestCtx {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("creator")

	policy := `
name: policy
resources:
- name: resource
`

	action := test.PolicySetupAction{
		Policy:        policy,
		PolicyCreator: "creator",
		ObjectsPerActor: map[string][]*types.Object{
			"bob": []*types.Object{
				types.NewObject("resource", "a"),
			},
		},
	}
	action.Run(ctx)

	return ctx
}

func TestTransfer_RootCannotTransferObject(t *testing.T) {
	ctx := transferSetup(t)

	ctx.SetRootPrincipal()
	alice := ctx.Actors.DID("alice")
	a := test.TransferObjectAction{
		PolicyId:    ctx.State.PolicyId,
		Object:      types.NewObject("resource", "a"),
		NewOwner:    alice,
		ExpectedErr: errors.ErrorType_UNAUTHORIZED,
	}
	a.Run(ctx)
}

func TestTransfer_NonOwnerCannotTransferObj(t *testing.T) {
	ctx := transferSetup(t)

	ctx.SetPrincipal("alice")
	alice := ctx.Actors.DID("alice")
	a := test.TransferObjectAction{
		PolicyId:    ctx.State.PolicyId,
		Object:      types.NewObject("resource", "a"),
		NewOwner:    alice,
		ExpectedErr: errors.ErrorType_UNAUTHORIZED,
	}
	a.Run(ctx)
}

func TestTransfer_OwnerCanTransfer(t *testing.T) {
	ctx := transferSetup(t)

	ctx.SetPrincipal("bob")
	alice := ctx.Actors.DID("alice")
	a := test.TransferObjectAction{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("resource", "a"),
		NewOwner: alice,
		Expected: &types.TransferObjectResponse{
			Record: &types.RelationshipRecord{
				PolicyId:     ctx.State.PolicyId,
				Relationship: types.NewActorRelationship("resource", "a", "owner", alice),
				Archived:     false,
				Metadata: &types.RecordMetadata{
					Creator: &types.Principal{
						Kind:       types.PrincipalKind_DID,
						Identifier: alice,
					},
					CreationTs: ctx.Time,
				},
			},
		},
	}
	a.Run(ctx)
}

func TestTransfer_ArchivedObjectCannotBeTransfered(t *testing.T) {
	ctx := transferSetup(t)
	ctx.SetPrincipal("bob")

	a1 := test.ArchiveObjectAction{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("resource", "a"),
	}
	a1.Run(ctx)

	alice := ctx.Actors.DID("alice")
	a := test.TransferObjectAction{
		PolicyId:    ctx.State.PolicyId,
		Object:      types.NewObject("resource", "a"),
		NewOwner:    alice,
		ExpectedErr: errors.ErrorType_OPERATION_FORBIDDEN,
	}
	a.Run(ctx)
}

func TestTransfer_UnregisteredObjectCannotBeTransfered(t *testing.T) {
	ctx := transferSetup(t)
	ctx.SetPrincipal("bob")

	alice := ctx.Actors.DID("alice")
	a := test.TransferObjectAction{
		PolicyId:    ctx.State.PolicyId,
		Object:      types.NewObject("resource", "unregistered"),
		NewOwner:    alice,
		ExpectedErr: errors.ErrorType_NOT_FOUND,
	}
	a.Run(ctx)
}
