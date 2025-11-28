package playground

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

func Test_SetAndVerify_ConsistentResults(t *testing.T) {
	ctx := test.NewTestCtx(t)

	data := types.SandboxData{
		PolicyDefinition: `
name: filesystem
resources:
- name: file
  permissions:
  - expr: reader
    name: read
  - name: write
  relations:
  - name: reader
    types:
    - actor
`,

		Relationships: `
file:readme#owner@did:user:bob // bob owns file readme
file:readme#reader@did:user:alice // alice can read and write file readme
`,
		PolicyTheorem: `
Authorizations {
  file:readme#reader@did:user:alice
}

Delegations {
}    
`,
	}

	resp, err := ctx.Playground.NewSandbox(ctx, &types.NewSandboxRequest{
		Name:        "test",
		Description: "",
	})
	require.NoError(t, err)
	handle := resp.Record.Handle

	setResp1, err := ctx.Playground.SetState(ctx, &types.SetStateRequest{
		Handle: handle,
		Data:   &data,
	})
	require.NoError(t, err)
	require.True(t, setResp1.Ok)

	verifyResult1, err := ctx.Playground.VerifyTheorems(ctx, &types.VerifyTheoremsRequest{
		Handle: handle,
	})
	require.NoError(t, err)
	require.True(t, verifyResult1.Result.Ok)

	setResp2, err := ctx.Playground.SetState(ctx, &types.SetStateRequest{
		Handle: handle,
		Data:   &data,
	})
	require.NoError(t, err)
	require.True(t, setResp2.Ok)

	r, _ := ctx.Playground.GetSandbox(ctx, &types.GetSandboxRequest{Handle: handle})
	t.Logf("%v", test.MustProtoToJson(r))

	verifyResult2, err := ctx.Playground.VerifyTheorems(ctx, &types.VerifyTheoremsRequest{
		Handle: handle,
	})
	require.NoError(t, err)
	//t.Logf("%v", test.MustProtoToJson(verifyResult2))
	require.True(t, verifyResult2.Result.Ok)

	require.Equal(t, verifyResult1, verifyResult2)

	setResp3, err := ctx.Playground.SetState(ctx, &types.SetStateRequest{
		Handle: handle,
		Data:   &data,
	})
	require.NoError(t, err)
	require.True(t, setResp3.Ok)

	verifyResult3, err := ctx.Playground.VerifyTheorems(ctx, &types.VerifyTheoremsRequest{
		Handle: handle,
	})
	require.NoError(t, err)

	require.Equal(t, verifyResult2, verifyResult3)
}
