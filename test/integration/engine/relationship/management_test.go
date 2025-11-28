package relationship

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
	"github.com/stretchr/testify/require"
)

func setupCheckManagement(t *testing.T) *test.TestCtx {
	ctx := test.NewTestCtx(t)
	admin := ctx.Actors.DID("admin")
	writer := ctx.Actors.DID("writer")
	reader := ctx.Actors.DID("reader")

	policy := `
name: policy
resources:
- name: file
  relations:
  - manages:
    - reader
    name: admin
    types:
    - actor
  - name: reader
    types:
    - actor
  - name: writer
    types:
    - actor
`

	action := test.PolicySetupAction{
		Policy:        policy,
		PolicyCreator: "creator",
		ObjectsPerActor: map[string][]*types.Object{
			"alice": {
				types.NewObject("file", "foo"),
			},
		},
		RelationshipsPerActor: map[string][]*types.Relationship{
			"alice": {
				types.NewActorRelationship("file", "foo", "reader", reader),
				types.NewActorRelationship("file", "foo", "writer", writer),
				types.NewActorRelationship("file", "foo", "admin", admin),
			},
		},
	}
	action.Run(ctx)

	return ctx
}

func Test_CheckManagementAuthority_ReturnsAuthorizedForAuthorizedActor(t *testing.T) {
	ctx := setupCheckManagement(t)

	result, err := ctx.Engine.CheckManagementAuthority(ctx, &types.CheckManagementAuthorityRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("file", "foo"),
		Relation: "reader",
		Actor:    types.NewActor(ctx.Actors.DID("admin")),
	})

	require.NoError(t, err)
	require.True(t, result.Authorized)
}

func Test_CheckManagementAuthority_ReturnsNotAuthorizedForUnauthorizedActor(t *testing.T) {
	ctx := setupCheckManagement(t)

	result, err := ctx.Engine.CheckManagementAuthority(ctx, &types.CheckManagementAuthorityRequest{
		PolicyId: ctx.State.PolicyId,
		Object:   types.NewObject("file", "foo"),
		Relation: "reader",
		Actor:    types.NewActor(ctx.Actors.DID("reader")),
	})

	require.NoError(t, err)
	require.False(t, result.Authorized)
}
