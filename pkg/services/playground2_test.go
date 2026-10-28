package services

import (
	"context"
	"net"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

func Test_Playground2_GenerateOverGRPC(t *testing.T) {
	if _, err := exec.LookPath("alloy"); err != nil {
		t.Skip("alloy not found in PATH")
	}
	listener := bufconn.Listen(1 << 20)
	server := grpc.NewServer()
	types.RegisterPlayground2Server(server, NewPlayground2Service(""))
	go server.Serve(listener)
	defer server.Stop()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return listener.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	resp, err := types.NewPlayground2Client(conn).Generate(context.Background(), &types.GenerateRequest{
		SystemModel: &types.SystemModel{
			Policy: `
name: notes
resources:
- name: note
  relations:
  - name: reader
    types: [actor]
  - name: editor
    types: [actor]
`,
			Constraints: []*types.SystemConstraint{{Constraint: &types.SystemConstraint_Disjoint{
				Disjoint: &types.DisjointConstraint{Resource: "note", Relations: []string{"reader", "editor"}},
			}}},
		},
		Config: &types.GenerateConfig{Count: 2},
	})
	require.NoError(t, err)
	require.Len(t, resp.Worlds, 2)
	for _, actor := range []string{"alice", "bob", "carol"} {
		require.Contains(t, resp.Worlds[0].Relationships, "@did:user:"+actor)
	}
}
