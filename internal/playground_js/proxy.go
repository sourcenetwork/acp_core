//go:build js

// package js
package playground_js

import (
	"context"
	"syscall/js"

	"github.com/cosmos/gogoproto/proto"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/services"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// PlaygroundConstructor returns a JS function which acts as a contructor for Playground Services.
// In JS land, the return of this constructor function is a JS object whose attributes implement
// the Playground protobuff definition.
// In other words, the returned type implements the TS generated interface for PlaygroundService.
// Note that, the JS playground contains an additional method, `close`, which should be called
// to free up resources in the Go runtime once the service is terminated.
//
// In Go land, the constructor function creates a PlaygroundServiceProxy object,
// which acts as proxy between the JS runtime and the Go code.
// The constructor returns the JS representation of the created PlaygroundServiceProxy
func PlaygroundConstructor(ctx context.Context) js.Func {
	return asyncFn(func(this js.Value, args []js.Value) (js.Value, error) {
		manager, err := runtime.NewRuntimeManager()
		if err != nil {
			return js.Undefined(), err
		}

		playground := NewPlaygroundServiceProxy(ctx, manager)
		return playground.GetJSValue(), nil
	})
}

// PlaygroundServiceProxy acts as a proxy between the JS runtime land and the underlying Go code
type PlaygroundServiceProxy struct {
	ctx      context.Context
	manager  runtime.RuntimeManager
	service  types.PlaygroundServiceServer
	proxyMap map[string]js.Func
	jsValue  js.Value
}

// NewPlaygroundServiceProxy creates a new PlaygroundService from a default context
func NewPlaygroundServiceProxy(ctx context.Context, manager runtime.RuntimeManager) *PlaygroundServiceProxy {
	service := services.NewPlaygroundService(manager)

	proxy := &PlaygroundServiceProxy{
		ctx:     ctx,
		manager: manager,
		service: service,
	}

	closeWrapper := js.FuncOf(func(this js.Value, args []js.Value) any {
		proxy.Close()
		return js.Undefined()
	})

	proxyMap := map[string]js.Func{
		"newSandbox":        asyncFn(wrapHandler(proxy.NewSandbox)),
		"listSandboxes":     asyncFn(wrapHandler(proxy.ListSandboxes)),
		"setState":          asyncFn(wrapHandler(proxy.SetState)),
		"restoreScratchpad": asyncFn(wrapHandler(proxy.RestoreScratchpad)),
		"getCatalogue":      asyncFn(wrapHandler(proxy.GetCatalogue)),
		"getSandbox":        asyncFn(wrapHandler(proxy.GetSandbox)),
		"verifyTheorems":    asyncFn(wrapHandler(proxy.VerifyTheorems)),
		"simulate":          asyncFn(wrapHandler(proxy.Simulate)),
		"close":             closeWrapper,
	}
	proxy.proxyMap = proxyMap
	proxy.makeJSValue()
	return proxy
}

// AsJSValue returns a JS Object whose attributes are js functions
// that dispatch execution to the playground methods.
func (s *PlaygroundServiceProxy) makeJSValue() {
	obj := make(map[string]any)
	for method, f := range s.proxyMap {
		obj[method] = f
	}
	s.jsValue = js.ValueOf(obj)
}

func (s *PlaygroundServiceProxy) GetJSValue() js.Value {
	return s.jsValue
}

func (s *PlaygroundServiceProxy) NewSandbox(this js.Value, args []js.Value) (*types.NewSandboxResponse, error) {
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

func (s *PlaygroundServiceProxy) ListSandboxes(this js.Value, args []js.Value) (*types.ListSandboxesResponse, error) {
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

func (s *PlaygroundServiceProxy) SetState(this js.Value, args []js.Value) (*types.SetStateResponse, error) {
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

func (s *PlaygroundServiceProxy) RestoreScratchpad(this js.Value, args []js.Value) (*types.RestoreScratchpadResponse, error) {
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

func (s *PlaygroundServiceProxy) GetCatalogue(this js.Value, args []js.Value) (*types.GetCatalogueResponse, error) {
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

func (s *PlaygroundServiceProxy) GetSandbox(this js.Value, args []js.Value) (*types.GetSandboxResponse, error) {
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

func (s *PlaygroundServiceProxy) VerifyTheorems(this js.Value, args []js.Value) (*types.VerifyTheoremsResponse, error) {
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

func (s *PlaygroundServiceProxy) Simulate(this js.Value, args []js.Value) (*types.SimulateReponse, error) {
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

// Close frees all resources used by the playground
func (s *PlaygroundServiceProxy) Close() {
	for _, f := range s.proxyMap {
		f.Release()
	}
	s.manager.Terminate()
}

// wrapHandler receives a handler method, which receives js values and return go values
// and turns into a hybrid which returns JS values or an error
func wrapHandler[R proto.Message](handler func(js.Value, []js.Value) (R, error)) func(js.Value, []js.Value) (js.Value, error) {
	return func(this js.Value, args []js.Value) (js.Value, error) {
		result, err := handler(this, args)
		if err != nil {
			return js.Undefined(), err
		}

		val, err := toJSObject(result)
		if err != nil {
			return js.Undefined(), err
		}

		return val, nil
	}
}
