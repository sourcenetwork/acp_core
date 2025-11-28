package object

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func setupQueryObjectOwner(t *testing.T) *test.TestCtx {
	ctx := test.NewTestCtx(t)

	ctx.SetPrincipal("admin")

	pol := `actor:
  doc: my actor
  name: actor-resource
description: ok
name: policy
resources:
- name: file
  permissions:
  - doc: own doc
    expr: owner
    name: own
  - expr: owner + reader
    name: read
  relations:
  - manages:
    - reader
    name: admin
  - doc: owner owns
    name: owner
    types:
    - actor-resource
  - name: reader
spec: none
`
	action := test.CreatePolicyAction{
		Policy: pol,
	}
	policy := action.Run(ctx)

	ctx.SetPrincipal("alice")
	objAction := test.RegisterObjectsAction{
		PolicyId: policy.Id,
		Objects: []*types.Object{
			types.NewObject("file", "1"),
		},
	}
	objAction.Run(ctx)

	return ctx
}

func TestGetObjectRegistration_ReturnsObjectOwner(t *testing.T) {
	ctx := setupQueryObjectOwner(t)

	resp, err := ctx.Engine.GetObjectRegistration(ctx, &types.GetObjectRegistrationRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("file", "1"),
	})

	want := &types.GetObjectRegistrationResponse{
		IsRegistered: true,
		OwnerId:      ctx.Actors.DID("alice"),
		Record: &types.RelationshipRecord{
			PolicyId:     ctx.State.PolicyId,
			Relationship: types.NewActorRelationship("file", "1", "owner", ctx.Actors.DID("alice")),
			Archived:     false,
			Metadata: &types.RecordMetadata{
				CreationTs: ctx.Time,
				Creator: &types.Principal{
					Identifier: ctx.Actors.DID("alice"),
					Kind:       types.PrincipalKind_DID,
				},
			},
		},
	}
	require.Equal(t, want, resp)
	require.NoError(t, err)
}

func TestGetObjectRegistration_QueryingForObjectInNonExistingResourceReturnsError(t *testing.T) {
	ctx := setupQueryObjectOwner(t)

	resp, err := ctx.Engine.GetObjectRegistration(ctx, &types.GetObjectRegistrationRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("missing-resource", "1"),
	})

	require.Nil(t, resp)
	require.NotNil(t, err)
}
func TestGetObjectOwner_QueryingPolicyThatDoesNotExistReturnError(t *testing.T) {
	ctx := setupQueryObjectOwner(t)

	resp, err := ctx.Engine.GetObjectRegistration(ctx, &types.GetObjectRegistrationRequest{
		PolicyId: "asbcf12345",
		Object:   types.NewObject("file", "1"),
	})

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_NOT_FOUND)
}

func TestGetObjectOwner_QueryingForUnregisteredObjectReturnsEmptyOwner(t *testing.T) {
	ctx := setupQueryObjectOwner(t)

	resp, err := ctx.Engine.GetObjectRegistration(ctx, &types.GetObjectRegistrationRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("file", "404"),
	})

	require.NoError(t, err)
	want := &types.GetObjectRegistrationResponse{
		IsRegistered: false,
		OwnerId:      "",
	}
	require.Equal(t, resp, want)
}
