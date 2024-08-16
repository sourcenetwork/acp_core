package simulator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/internal/theorem"
	"github.com/sourcenetwork/acp_core/pkg/playground"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func simulate(t *testing.T, req *playground.SandboxData) *types.AnnotatedPolicyTheoremResult {
	ctx := context.TODO()
	manager, err := runtime.NewRuntimeManager()
	require.NoError(t, err)

	resp1, err := HandleNewSandboxRequest(ctx, manager, &playground.NewSandboxRequest{Name: "test"})
	require.NoError(t, err)

	handler := SetStateHandler{}
	resp2, err := handler.Handle(ctx, manager, &playground.SetStateRequest{
		Data:   req,
		Handle: resp1.Record.Handle,
	})
	require.NoError(t, err)
	_ = resp2

	resp3, err := HandleVerifyTheorem(ctx, manager, &playground.VerifyTheoremsRequest{
		Handle: resp1.Record.Handle,
	})
	require.NoError(t, err)

	return resp3.Result
}

func TestSimulator_Success(t *testing.T) {
	req := playground.SandboxData{
		PolicyDefinition: `
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
		Relationships: `
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
	}

	result := simulate(t, &req)
	report := theorem.GenerateTheoremReport(types.FromAnnotatedResult(result))
	t.Logf("%v", report)
	t.Fail()
}
