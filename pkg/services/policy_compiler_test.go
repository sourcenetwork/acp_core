package services

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

func Test_PolicyCompiler_CompileOverGRPC(t *testing.T) {
	listener := bufconn.Listen(1 << 20)
	server := grpc.NewServer()
	types.RegisterPolicyCompilerServer(server, NewPolicyCompilerService())
	go server.Serve(listener)
	defer server.Stop()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return listener.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	resp, err := types.NewPolicyCompilerClient(conn).Compile(context.Background(), &types.CompileRequest{
		Policy: `
name: notes
resources:
- name: note
  relations:
  - name: reader
    types: [actor]
  permissions:
  - name: read
    expr: reader
`,
	})
	require.NoError(t, err)
	require.Len(t, resp.Files, 2)
	require.Equal(t, "notes.als", resp.Files[0].Name)
	require.Contains(t, resp.Files[0].Content, "fun note_read: Note->Actor")
	require.Equal(t, "notes_management.als", resp.Files[1].Name)
}

func Test_PolicyCompiler_CompileEmitsAssumptions(t *testing.T) {
	policy := `
name: notes
resources:
- name: note
  relations:
  - name: reader
    types: [actor]
  - name: editor
    types: [actor]
`
	resp, err := NewPolicyCompilerService().Compile(context.Background(), &types.CompileRequest{
		Policy: policy,
		Assumptions: []*types.Assumption{{Assumption: &types.Assumption_Disjoint{
			Disjoint: &types.DisjointAssumption{Resource: "note", Relations: []string{"reader", "editor"}},
		}}},
	})
	require.NoError(t, err)
	require.Len(t, resp.Files, 3)
	require.Equal(t, "notes_assumptions.als", resp.Files[2].Name)
	require.Contains(t, resp.Files[2].Content, "no (Note <: reader) & (Note <: editor)")

	_, err = NewPolicyCompilerService().Compile(context.Background(), &types.CompileRequest{
		Policy: policy,
		Assumptions: []*types.Assumption{{Assumption: &types.Assumption_Acyclic{
			Acyclic: &types.AcyclicAssumption{Resource: "note", Relation: "reader"},
		}}},
	})
	require.ErrorContains(t, err, "can't form cycles")
}
