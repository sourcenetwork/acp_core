package simulator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

func simulate(ctx context.Context, req *types.SimulateRequest) (*types.AnnotatedSimulationResult, error) {
	resp, err := HandleSimulateRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func TestSimulator_Success(t *testing.T) {
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

	result, err := simulate(context.TODO(), &req)
	require.NoError(t, err)
	t.Logf("%v", GenerateSimulationReport(result))
	t.Fail()
}

/*
func TestSimulator_InvalidGrammars_Policy_Rels_Theorems_ReturnErrors(t *testing.T) {
	req := types.SimulateRequest{
		Declaration: &types.SimulationCtxDeclaration{
			Policy: `
        name: test
        resources:
          file
        `,
			MarshalType: types.PolicyMarshalingType_SHORT_YAML,
			RelationshipSet: `
		file:abc#relationship@did:ex:bob
		file:abc#reader@did:ex:alice
		`,
			PolicyTheorem: `
		Authorizations {
		  file:abc#read@invalid-actor
		}
		Delegations {
		  did:ex:bob < file:abc#reader
		}
		`,
		},
	}

	resp, err := simulate(context.TODO(), &req)
	require.NoError(t, err)
	respJson, err := json.MarshalIndent(resp, "", "  ")
	require.NoError(t, err)
	t.Logf("%v", string(respJson))
}

*/
