package object

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
	testutil "github.com/sourcenetwork/acp_core/test/util"
)

func revealSetup(t *testing.T) *test.TestCtx {
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

func TestReveal_RevealRegistrationSuceedsWithProvidedTs(t *testing.T) {
	ctx := transferSetup(t)

	alice := ctx.SetPrincipal("alice")
	a := test.RevealRegistrationAction{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("resource", "b"),
		Ts:       testutil.MustDateTimeToProto("2025-01-01 00:00:00"),
		Expected: &types.RelationshipRecord{
			PolicyId: ctx.State.PolicyId,
			Archived: false,
			Metadata: &types.RecordMetadata{
				Creator:      &alice,
				CreationTs:   testutil.MustDateTimeToProto("2025-01-01 00:00:00"),
				LastModified: ctx.Time,
				Supplied:     nil,
			},
			Relationship: types.NewActorRelationship("resource", "b", "owner", alice.Identifier),
		},
	}
	a.Run(ctx)
}

func TestReveal_RootCannotRevealRegistration(t *testing.T) {
	ctx := transferSetup(t)

	ctx.SetRootPrincipal()
	a := test.RevealRegistrationAction{
		PolicyId:    ctx.State.PolicyId,
		Object:      types.NewObject("resource", "b"),
		Ts:          test.DefaultTs,
		ExpectedErr: errors.ErrorType_UNAUTHORIZED,
	}
	a.Run(ctx)
}

func TestReveal_ErrorsIfObjectIsAlreadyRegistered(t *testing.T) {
	ctx := transferSetup(t)

	ctx.SetPrincipal("alice")
	a := test.RevealRegistrationAction{
		PolicyId:    ctx.State.PolicyId,
		Object:      types.NewObject("resource", "a"),
		Ts:          test.DefaultTs,
		ExpectedErr: errors.ErrorType_OPERATION_FORBIDDEN,
	}
	a.Run(ctx)
}

func TestReveal_ArchivedObjectCannotBeRevealed(t *testing.T) {
	ctx := transferSetup(t)
	ctx.SetPrincipal("bob")

	a1 := test.ArchiveObjectAction{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("resource", "a"),
	}
	a1.Run(ctx)

	a := test.RevealRegistrationAction{
		PolicyId:    ctx.State.PolicyId,
		Object:      types.NewObject("resource", "a"),
		Ts:          test.DefaultTs,
		ExpectedErr: errors.ErrorType_OPERATION_FORBIDDEN,
	}
	a.Run(ctx)
}
