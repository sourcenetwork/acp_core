package playground

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func Test_Simulate_ReturnsErrorIfSubjectInOwnerRelationshipIsNotActor(t *testing.T) {
	ctx := test.NewTestCtx(t)

	data := types.SandboxData{
		PolicyDefinition: `
name: shinzo
resources:
  file:
    relations:
      owner:
		types:
        - actor
  group:	
    relations:
      owner:
        types:
        - actor
`,
		Relationships: `
file:logs#owner@group:example#owner
`,
		PolicyTheorem: `
		Authorizations {
		}
		Delegations {
		}
		`,
	}

	resp, err := ctx.Playground.Simulate(ctx, &types.SimulateRequest{
		Data: &data,
	})
	require.NoError(t, err)
	require.False(t, resp.ValidData)
	require.Contains(t, resp.Errors.RelationshipsErrors[0].Message, "owner relationship requires a `did` actor")
}
