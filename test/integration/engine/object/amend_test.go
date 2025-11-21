package object

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func amendSetup(t *testing.T) *test.TestCtx {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("creator")

	policy := `name: policy
resources:
- name: resource
  relations:
  - name: owner
    types:
    - actor
spec: none
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

func TestAmend_RootCanAmendRegistration(t *testing.T) {
	ctx := amendSetup(t)

	ctx.SetRootPrincipal()
	alice := ctx.Actors.DID("alice")
	a := test.AmendRegistrationAction{
		PolicyId:     ctx.State.PolicyId,
		Object:       types.NewObject("resource", "a"),
		NewOwner:     alice,
		NewTimestamp: test.DefaultTs,
		Expected: &types.AmendRegistrationResponse{
			Record: &types.RelationshipRecord{
				PolicyId:     ctx.State.PolicyId,
				Relationship: types.NewActorRelationship("resource", "a", "owner", alice),
				Archived:     false,
				Metadata: &types.RecordMetadata{
					CreationTs: test.DefaultTs,
					Creator: &types.Principal{
						Identifier: alice,
						Kind:       types.PrincipalKind_DID,
					},
					LastModified: ctx.Time,
				},
			},
		},
	}
	a.Run(ctx)
}

func TestAmend_NonOwnerCannotAmendRegistration(t *testing.T) {
	ctx := amendSetup(t)

	ctx.SetPrincipal("alice")
	alice := ctx.Actors.DID("alice")
	a := test.AmendRegistrationAction{
		PolicyId:     ctx.State.PolicyId,
		Object:       types.NewObject("resource", "a"),
		NewOwner:     alice,
		NewTimestamp: test.DefaultTs,
		ExpectedErr:  errors.ErrorType_UNAUTHORIZED,
	}
	a.Run(ctx)
}

func TestAmend_OwnerCannotAmendRegistration(t *testing.T) {
	ctx := amendSetup(t)

	ctx.SetPrincipal("bob")
	alice := ctx.Actors.DID("alice")
	a := test.AmendRegistrationAction{
		PolicyId:     ctx.State.PolicyId,
		Object:       types.NewObject("resource", "a"),
		NewOwner:     alice,
		NewTimestamp: test.DefaultTs,
		ExpectedErr:  errors.ErrorType_UNAUTHORIZED,
	}
	a.Run(ctx)
}

func TestAmend_ArchivedObjectCannotBeAmended(t *testing.T) {
	ctx := amendSetup(t)
	ctx.SetPrincipal("bob")

	a1 := test.ArchiveObjectAction{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("resource", "a"),
	}
	a1.Run(ctx)

	ctx.SetRootPrincipal()
	alice := ctx.Actors.DID("alice")
	a := test.AmendRegistrationAction{
		PolicyId:     ctx.State.PolicyId,
		Object:       types.NewObject("resource", "a"),
		NewOwner:     alice,
		NewTimestamp: test.DefaultTs,
		ExpectedErr:  errors.ErrorType_OPERATION_FORBIDDEN,
	}
	a.Run(ctx)
}

func TestAmend_UnregisteredObjectCannotBeAmended(t *testing.T) {
	ctx := amendSetup(t)
	ctx.SetRootPrincipal()

	alice := ctx.Actors.DID("alice")
	a := test.AmendRegistrationAction{
		PolicyId:     ctx.State.PolicyId,
		Object:       types.NewObject("resource", "unregistered"),
		NewOwner:     alice,
		NewTimestamp: test.DefaultTs,
		ExpectedErr:  errors.ErrorType_BAD_INPUT,
	}
	a.Run(ctx)
}

func TestAmend_OwnerCannotAmendObjectToThemselves(t *testing.T) {
	ctx := amendSetup(t)
	ctx.SetPrincipal("bob")
	bob := ctx.Actors.DID("bob")

	a := test.AmendRegistrationAction{
		PolicyId:     ctx.State.PolicyId,
		Object:       types.NewObject("resource", "a"),
		NewOwner:     bob,
		NewTimestamp: test.DefaultTs,
		ExpectedErr:  errors.ErrorType_UNAUTHORIZED,
	}
	a.Run(ctx)
}

func TestAmend_ErrorsWithoutNewTimestamp(t *testing.T) {
	ctx := amendSetup(t)
	ctx.SetRootPrincipal()
	bob := ctx.Actors.DID("bob")

	a := test.AmendRegistrationAction{
		PolicyId:    ctx.State.PolicyId,
		Object:      types.NewObject("resource", "a"),
		NewOwner:    bob,
		ExpectedErr: errors.ErrorType_BAD_INPUT,
	}
	a.Run(ctx)
}
