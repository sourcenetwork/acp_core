//go:build js

// package js
package js

import (
	"context"
	"syscall/js"

	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/services"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// NewPlayground returns a JS function which acts as a contructor for playgrounds.
// In JS land, the return of this constructor function is a JS object whose attributes
// are named similarly to the Playgroung protobuff definition.
// Calling these attributes will execute the expected Playground sevice operation.
//
// In Go land, the constructor function creates a PlaygroundServiceProxy object,
// which acts as proxy between the JS runtime and the Go code.
// The constructor returns the JS representation of the created PlaygroundServiceProxy
func NewPlayground(ctx context.Context) js.Func {
	return asyncFn(func(this js.Value, args []js.Value) (any, error) {
		playground, err := newPlaygroundServiceProxy(ctx)
		if err != nil {
			return nil, err
		}
		return playground.asValue(), nil
	})
}

// PlaygroundServiceProxy acts as a proxy between the JS runtime land and the underlying Go code
type PlaygroundServiceProxy struct {
	ctx      context.Context
	manager  runtime.RuntimeManager
	service  types.PlaygroundServiceServer
	proxyMap map[string]js.Func
}

// newPlaygroundServiceProxy creates a new PlaygroundService from a default context
func newPlaygroundServiceProxy(ctx context.Context) (*PlaygroundServiceProxy, error) {
	manager, err := runtime.NewRuntimeManager()
	if err != nil {
		return nil, err
	}
	service := services.NewPlaygroundService(manager)

	proxy := &PlaygroundServiceProxy{
		ctx:     ctx,
		manager: manager,
		service: service,
	}
	proxyMap := map[string]js.Func{
		"newSandbox":        asyncHandler(proxy.newSandbox),
		"listSandboxes":     asyncHandler(proxy.listSandboxes),
		"setState":          asyncHandler(proxy.setState),
		"restoreScratchpad": asyncHandler(proxy.restoreScratchpad),
		"getCatalogue":      asyncHandler(proxy.getCatalogue),
		"getSandbox":        asyncHandler(proxy.getSandbox),
		"verifyTheorems":    asyncHandler(proxy.verifyTheorems),
		"simulate":          asyncHandler(proxy.simulate),
		"close":             proxy.close(),
	}
	proxy.proxyMap = proxyMap
	return proxy, nil
}

// asValue returns a JS Object whose attributes are js functions
// that dispatch execution to the playground methods.
func (s *PlaygroundServiceProxy) asValue() js.Value {
	obj := make(map[string]any)
	for method, f := range s.proxyMap {
		obj[method] = f
	}
	return js.ValueOf(obj)
}

func (s *PlaygroundServiceProxy) newSandbox(this js.Value, args []js.Value) (*types.NewSandboxResponse, error) {
	req := &types.NewSandboxRequest{}
	err := unmarsahlArgs(req, args)
	if err != nil {
		return nil, err
	}

	resp, err := s.service.NewSandbox(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *PlaygroundServiceProxy) listSandboxes(this js.Value, args []js.Value) (*types.ListSandboxesResponse, error) {
	req := &types.ListSandboxesRequest{}
	err := unmarsahlArgs(req, args)
	if err != nil {
		return nil, err
	}

	resp, err := s.service.ListSandboxes(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *PlaygroundServiceProxy) setState(this js.Value, args []js.Value) (*types.SetStateResponse, error) {
	req := &types.SetStateRequest{}
	err := unmarsahlArgs(req, args)

	if err != nil {
		return nil, err
	}

	resp, err := s.service.SetState(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *PlaygroundServiceProxy) restoreScratchpad(this js.Value, args []js.Value) (*types.RestoreScratchpadResponse, error) {
	req := &types.RestoreScratchpadRequest{}
	err := unmarsahlArgs(req, args)
	if err != nil {
		return nil, err
	}

	resp, err := s.service.RestoreScratchpad(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *PlaygroundServiceProxy) getCatalogue(this js.Value, args []js.Value) (*types.GetCatalogueResponse, error) {
	req := &types.GetCatalogueRequest{}
	err := unmarsahlArgs(req, args)
	if err != nil {
		return nil, err
	}

	resp, err := s.service.GetCatalogue(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *PlaygroundServiceProxy) getSandbox(this js.Value, args []js.Value) (*types.GetSandboxResponse, error) {
	req := &types.GetSandboxRequest{}
	err := unmarsahlArgs(req, args)
	if err != nil {
		return nil, err
	}

	resp, err := s.service.GetSandbox(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *PlaygroundServiceProxy) verifyTheorems(this js.Value, args []js.Value) (*types.VerifyTheoremsResponse, error) {
	req := &types.VerifyTheoremsRequest{}
	err := unmarsahlArgs(req, args)
	if err != nil {
		return nil, err
	}

	resp, err := s.service.VerifyTheorems(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *PlaygroundServiceProxy) simulate(this js.Value, args []js.Value) (*types.SimulateReponse, error) {
	req := &types.SimulateRequest{}
	err := unmarsahlArgs(req, args)
	if err != nil {
		return nil, err
	}

	resp, err := s.service.Simulate(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// close frees all resources used by the playground
func (s *PlaygroundServiceProxy) close() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		for _, f := range s.proxyMap {
			f.Release()
		}
		return nil
	})
}
