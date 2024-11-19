//go:build js

package js

import (
	"context"
	"syscall/js"
	"testing"

	"github.com/cosmos/gogoproto/jsonpb"
	"github.com/cosmos/gogoproto/proto"

	"github.com/sourcenetwork/acp_core/internal/playground_js"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ types.PlaygroundServiceServer = (*PlaygroundJS)(nil)

// NewPlaygruondJS returns a new playground JS Wrapper
func NewPlaygroundJS(t *testing.T, manager runtime.RuntimeManager) *PlaygroundJS {
	proxy := playground_js.NewPlaygroundServiceProxy(context.TODO(), manager)

	return &PlaygroundJS{
		proxy: proxy,
		this:  proxy.GetJSValue(),
	}
}

// PlaygroundJS is a wrapper of the js wrapper
// It implements the Go definition of a Playground service.
// The wrapper receives Go objects, turns them into JS objects
// and calls the JS implementation of the Playground
type PlaygroundJS struct {
	proxy *playground_js.PlaygroundServiceProxy
	this  js.Value
}

func (s *PlaygroundJS) NewSandbox(ctx context.Context, req *types.NewSandboxRequest) (*types.NewSandboxResponse, error) {
	return s.proxy.NewSandbox(s.this, mustMapArgument(req))
}

func (s *PlaygroundJS) ListSandboxes(ctx context.Context, req *types.ListSandboxesRequest) (*types.ListSandboxesResponse, error) {
	return s.proxy.ListSandboxes(s.this, mustMapArgument(req))
}

func (s *PlaygroundJS) SetState(ctx context.Context, req *types.SetStateRequest) (*types.SetStateResponse, error) {
	return s.proxy.SetState(s.this, mustMapArgument(req))
}

func (s *PlaygroundJS) GetCatalogue(ctx context.Context, req *types.GetCatalogueRequest) (*types.GetCatalogueResponse, error) {
	return s.proxy.GetCatalogue(s.this, mustMapArgument(req))
}

func (s *PlaygroundJS) VerifyTheorems(ctx context.Context, req *types.VerifyTheoremsRequest) (*types.VerifyTheoremsResponse, error) {
	return s.proxy.VerifyTheorems(s.this, mustMapArgument(req))
}

func (s *PlaygroundJS) RestoreScratchpad(ctx context.Context, req *types.RestoreScratchpadRequest) (*types.RestoreScratchpadResponse, error) {
	return s.proxy.RestoreScratchpad(s.this, mustMapArgument(req))
}

func (s *PlaygroundJS) Simulate(ctx context.Context, req *types.SimulateRequest) (*types.SimulateReponse, error) {
	return s.proxy.Simulate(s.this, mustMapArgument(req))
}

func (s *PlaygroundJS) GetSandbox(ctx context.Context, req *types.GetSandboxRequest) (*types.GetSandboxResponse, error) {
	return s.proxy.GetSandbox(s.this, mustMapArgument(req))
}

func (s *PlaygroundJS) GetSampleSandboxes(ctx context.Context, req *types.GetSampleSandboxesRequest) (*types.GetSampleSandboxesResponse, error) {
	return s.proxy.GetSampleSandboxes(s.this, mustMapArgument(req))
}

func mustMapArgument(req proto.Message) []js.Value {
	marshaler := jsonpb.Marshaler{}
	valStr, err := marshaler.MarshalToString(req)
	if err != nil {
		panic(err)
	}

	jsVal := js.Global().Get("JSON").Call("parse", valStr)
	return []js.Value{jsVal}
}
