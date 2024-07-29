package simulator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func TestSimulator_Success(t *testing.T) {
	ctx := test.NewTestCtx(t)

	req := types.SimulateRequest{
		Declaration: &types.SimulationCtxDeclaration{
			Policy: `
        name: test
        resources:
          file:
            relations:
              owner:
                types:
                  - actor
              reader:
                types:
                  - actor
            permissions:
              read:
                expr: owner + reader
              write:
                expr: owner
        `,
			MarshalType: types.PolicyMarshalingType_SHORT_YAML,
			RelationshipSet: `
		file:abc#owner@did:ex:bob
		file:abc#reader@did:ex:alice
		`,
			PolicyTheorem: `
		Authorizations {
		  file:abc#read@did:ex:bob
		  file:abc#write@did:ex:bob
		  file:abc#read@did:ex:alice
		  file:abc#unknown_relation@did:ex:alice
		  ! file:abc#write@did:ex:alice
		}
		Delegations {
		  did:ex:bob > file:abc#reader
		  did:ex:alice > file:abc#reader
		}
		`,
		},
	}

	resp, err := ctx.Engine.Simulate(ctx, &req)
	require.NoError(t, err)
	respJson, err := json.MarshalIndent(resp, "", "  ")
	require.NoError(t, err)
	t.Logf("%v", string(respJson))
	t.Fail()
	//want := &types.SimulateResponse{ }
	//require.Equal(t, want, resp)
}
