package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/internal/sandbox"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func evaluate(t *testing.T, req *types.EvaluateRequest) *types.EvaluateResponse {
	t.Helper()
	manager, err := runtime.NewRuntimeManager()
	require.NoError(t, err)
	resp, err := NewSimulatorService(manager).Evaluate(context.Background(), req)
	require.NoError(t, err)
	return resp
}

func sampleRequest() *types.EvaluateRequest {
	data := sandbox.Samples[0].Data
	return &types.EvaluateRequest{
		Policy:        data.PolicyDefinition,
		Relationships: data.Relationships,
		Theorems:      data.PolicyTheorem,
	}
}

func Test_Simulator_EvaluatesTheorems(t *testing.T) {
	resp := evaluate(t, sampleRequest())

	require.True(t, resp.Valid)
	require.True(t, resp.Result.Ok)
	require.NotZero(t, resp.Result.TheoremCount)
	require.NotNil(t, resp.Ctx.Policy)
}

func Test_Simulator_ReportsFailedTheorems(t *testing.T) {
	req := sampleRequest()
	req.Theorems = `Authorizations {
  file:readme#write@did:user:eve
}
Delegations {}`

	resp := evaluate(t, req)

	require.True(t, resp.Valid)
	require.False(t, resp.Result.Ok)
	require.Equal(t, uint32(1), resp.Result.Failures)
}

func Test_Simulator_EmptyTheoremsEvaluatesNothing(t *testing.T) {
	req := sampleRequest()
	req.Theorems = ""

	resp := evaluate(t, req)

	require.True(t, resp.Valid)
	require.Zero(t, resp.Result.TheoremCount)
}

func Test_Simulator_InvalidRelationshipsAreReported(t *testing.T) {
	req := sampleRequest()
	req.Relationships = "file:readme#unknown@did:user:bob"

	resp := evaluate(t, req)

	require.False(t, resp.Valid)
	require.NotEmpty(t, resp.Errors.RelationshipsErrors)
	require.Nil(t, resp.Result)
}

func Test_Simulator_IsDeterministic(t *testing.T) {
	first, err := evaluate(t, sampleRequest()).Marshal()
	require.NoError(t, err)
	second, err := evaluate(t, sampleRequest()).Marshal()
	require.NoError(t, err)

	require.Equal(t, first, second)
}
