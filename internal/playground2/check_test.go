package playground2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/internal/telemetry"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func folderScenario(relationships string) *types.Scenario {
	return &types.Scenario{
		SystemModel: &types.SystemModel{Policy: folderPolicy},
		World:       &types.World{Relationships: relationships},
	}
}

const folderWorld = `folder:specs#owner@did:user:alice
folder:specs#viewer@did:user:bob
doc:roadmap#owner@did:user:carol
doc:roadmap#folder@folder:specs
`

func requestCtx() context.Context {
	return telemetry.InitRequestContext(context.Background())
}

func Test_Check(t *testing.T) {
	cases := []struct {
		actor, object, permission string
		allowed                   bool
	}{
		{"did:user:bob", "roadmap", "view", true},    // through the folder
		{"did:user:carol", "roadmap", "view", true},  // owner
		{"did:user:dave", "roadmap", "view", false},  // nothing
		{"did:user:bob", "roadmap", "viewer", false}, // relations can be checked: bob isn't a direct viewer
	}
	for _, c := range cases {
		resp, err := Check(requestCtx(), &types.CheckRequest{
			Scenario:   folderScenario(folderWorld),
			Object:     types.NewObject("doc", c.object),
			Permission: c.permission,
			Actor:      types.NewActor(c.actor),
		})
		require.NoError(t, err)
		require.Equal(t, c.allowed, resp.Allowed, "%v %v %v", c.actor, c.permission, c.object)
	}
}

func Test_Explain(t *testing.T) {
	resp, err := Explain(requestCtx(), &types.ExplainRequest{
		Scenario:   folderScenario(folderWorld),
		Object:     types.NewObject("doc", "roadmap"),
		Permission: "view",
		Actor:      types.NewActor("did:user:bob"),
	})

	require.NoError(t, err)
	require.True(t, resp.Allowed)
	require.NotEmpty(t, resp.Graph.RootNodeId)
	require.NotEmpty(t, resp.Graph.Nodes)
	require.NotEmpty(t, resp.Graph.Edges)
}

func Test_Check_InvalidScenario(t *testing.T) {
	_, err := Check(requestCtx(), &types.CheckRequest{
		Scenario:   folderScenario("doc:roadmap#viewer@did:user:bob\n"),
		Object:     types.NewObject("doc", "roadmap"),
		Permission: "view",
		Actor:      types.NewActor("did:user:bob"),
	})
	require.ErrorIs(t, err, ErrScenario)
	require.ErrorContains(t, err, "relationships 1:0: object not registered")

	_, err = Check(requestCtx(), &types.CheckRequest{Object: types.NewObject("doc", "roadmap")})
	require.ErrorIs(t, err, ErrScenario)
}
