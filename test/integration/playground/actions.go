package playground

import (
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
	"github.com/stretchr/testify/require"
)

type NewSandbox struct {
	Req         *types.NewSandboxRequest
	Expected    *types.NewSandboxResponse
	ExpectedErr error
}

func (a *NewSandbox) Run(ctx *test.TestCtx) *types.NewSandboxResponse {
	resp, err := ctx.Playground.NewSandbox(ctx, a.Req)
	test.AssertResults(ctx, resp, a.Expected, err, a.ExpectedErr)
	return resp
}

type ListSandboxes struct {
	Req         *types.ListSandboxesRequest
	ExpectedErr error
	ExpectedLen int
}

func (a *ListSandboxes) Run(ctx *test.TestCtx) *types.ListSandboxesResponse {
	resp, err := ctx.Playground.ListSandboxes(ctx, a.Req)
	test.AssertResults(ctx, len(resp.Records), a.ExpectedLen, err, a.ExpectedErr)
	return resp
}

type SetState struct {
	Req         *types.SetStateRequest
	ExpectedErr error
	Assertions  []Assertion
}

func (a *SetState) Run(ctx *test.TestCtx) *types.SetStateResponse {
	resp, err := ctx.Playground.SetState(ctx, a.Req)
	if a.ExpectedErr != nil {
		require.ErrorIs(ctx.T, err, a.ExpectedErr)
		return nil
	} else {
		require.NoError(ctx.T, err)
	}

	if a.Assertions == nil || len(a.Assertions) == 0 {
		require.True(ctx.T, resp.Ok, "expected ok response: got %v", resp)
		require.Equal(ctx.T, &types.SandboxDataErrors{}, resp.Errors)
		return resp
	}

	for _, assertion := range a.Assertions {
		assertion(ctx.T, resp.Errors)
	}

	return resp
}

type GetCatalogue struct {
	Ctx         *test.TestCtx
	Req         *types.GetCatalogueRequest
	Expected    *types.GetCatalogueResponse
	ExpectedErr error
}

func (a *GetCatalogue) Run(ctx *test.TestCtx) *types.GetCatalogueResponse {
	resp, err := ctx.Playground.GetCatalogue(ctx, a.Req)
	test.AssertResults(ctx, resp, a.Expected, err, a.ExpectedErr)
	return resp
}

type VerifyTheorems struct {
	Req         *types.VerifyTheoremsRequest
	Expected    *types.VerifyTheoremsResponse
	ExpectedErr error
}

func (a *VerifyTheorems) Run(ctx *test.TestCtx) *types.VerifyTheoremsResponse {
	resp, err := ctx.Playground.VerifyTheorems(ctx, a.Req)
	test.AssertResults(ctx, resp, a.Expected, err, a.ExpectedErr)
	return resp
}

type NewAndSet struct {
	Data       *types.SandboxData
	Assertions []Assertion
}

func (a *NewAndSet) Run(ctx *test.TestCtx) uint64 {
	a1 := NewSandbox{
		Req: &types.NewSandboxRequest{},
	}
	resp := a1.Run(ctx)

	a2 := SetState{
		Req: &types.SetStateRequest{
			Handle: resp.Record.Handle,
			Data:   a.Data,
		},
		Assertions: a.Assertions,
	}
	a2.Run(ctx)

	return resp.Record.Handle
}
