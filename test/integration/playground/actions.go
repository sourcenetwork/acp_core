package playground

import (
	"github.com/sourcenetwork/acp_core/pkg/playground"
	"github.com/sourcenetwork/acp_core/test"
)

type NewSandbox struct {
	Req         *playground.NewSandboxRequest
	Expected    *playground.NewSandboxResponse
	ExpectedErr error
}

func (a *NewSandbox) Run(ctx *test.TestCtx) *playground.NewSandboxResponse {
	resp, err := ctx.Playground.NewSandbox(ctx, a.Req)
	test.AssertResults(ctx, resp, a.Req, err, a.ExpectedErr)
	return resp
}

type ListSandboxes struct {
	Req         *playground.ListSandboxesRequest
	Expected    *playground.ListSandboxesResponse
	ExpectedErr error
}

func (a *ListSandboxes) Run(ctx *test.TestCtx) *playground.ListSandboxesResponse {
	resp, err := ctx.Playground.ListSandboxes(ctx, a.Req)
	test.AssertResults(ctx, resp, a.Req, err, a.ExpectedErr)
	return resp
}

type SetState struct {
	Req         *playground.SetStateRequest
	Expected    *playground.SetStateResponse
	ExpectedErr error
}

func (a *SetState) Run(ctx *test.TestCtx) *playground.SetStateResponse {
	resp, err := ctx.Playground.SetState(ctx, a.Req)
	test.AssertResults(ctx, resp, a.Req, err, a.ExpectedErr)
	return resp
}

type GetCatalogue struct {
	Ctx         *test.TestCtx
	Req         *playground.GetCatalogueRequest
	Expected    *playground.GetCatalogueResponse
	ExpectedErr error
}

func (a *GetCatalogue) Run(ctx *test.TestCtx) *playground.GetCatalogueResponse {
	resp, err := ctx.Playground.GetCatalogue(ctx, a.Req)
	test.AssertResults(ctx, resp, a.Req, err, a.ExpectedErr)
	return resp
}

type VerifyTheorems struct {
	Req         *playground.VerifyTheoremsRequest
	Expected    *playground.VerifyTheoremsResponse
	ExpectedErr error
}

func (a *VerifyTheorems) Run(ctx *test.TestCtx) *playground.VerifyTheoremsResponse {
	resp, err := ctx.Playground.VerifyTheorems(ctx, a.Req)
	test.AssertResults(ctx, resp, a.Req, err, a.ExpectedErr)
	return resp
}

type NewAndSet struct {
	Data *playground.SandboxData
}

func (a *NewAndSet) Run(ctx *test.TestCtx) uint64 {
	a1 := NewSandbox{
		Req: &playground.NewSandboxRequest{},
	}
	resp := a1.Run(ctx)

	a2 := SetState{
		Req: &playground.SetStateRequest{
			Handle: resp.Record.Handle,
			Data:   a.Data,
		},
	}
	a2.Run(ctx)

	return resp.Record.Handle
}
