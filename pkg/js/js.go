//go:build js

// package js
package js

import (
	"context"
	"syscall/js"

	"github.com/cosmos/gogoproto/jsonpb"
	"github.com/cosmos/gogoproto/proto"

	"github.com/sourcenetwork/acp_core/pkg/errors"
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

func (s *PlaygroundServiceProxy) close() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		for _, f := range s.proxyMap {
			f.Release()
		}
		return nil
	})
}

func toJSObject[T proto.Message](val T) (js.Value, error) {
	marshaler := jsonpb.Marshaler{}
	valStr, err := marshaler.MarshalToString(val)
	if err != nil {
		return js.Value{}, err
	}
	jsVal := js.Global().Get("JSON").Call("parse", valStr)
	return jsVal, nil
}

func unmarsahlArgs(container proto.Message, args []js.Value) error {
	count := len(args)
	if count != 1 {
		return newInvalidArgsErr(count)
	}

	arg := args[0]

	reqJson := js.Global().Get("JSON").Call("stringify", arg)
	if reqJson.Type() != js.TypeString {
		panic("expected string")
	}

	err := jsonpb.UnmarshalString(reqJson.String(), container)
	if err != nil {
		return errors.NewFromCause("could not unmarshal serialized req obj", err, errors.ErrorType_BAD_INPUT)
	}

	return nil
}

// asyncHandler takes a handler and turns it into a JS function which returns a promise.
// First a Promise executor function is created, inside of which a worker go routine is created and dispatched.
// The executor returns immediately and the returned Promise is returned as the result of the function
//
// Inside the worker go routine, the Handler is invoked and according to the results
// either invokes the reject if the error wasn't nil or resolves the promise with the result
func asyncHandler[R proto.Message](handler func(js.Value, []js.Value) (R, error)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		var promiseExecutor js.Func
		promiseExecutor = js.FuncOf(func(_ js.Value, promiseArgs []js.Value) any {
			go func() {
				defer promiseExecutor.Release()
				resolve := promiseArgs[0]
				reject := promiseArgs[1]

				result, err := handler(this, args)
				if err != nil {
					jsErr := newJSError(err)
					reject.Invoke(jsErr)
					return
				}

				val, err := toJSObject(result)
				if err != nil {
					jsErr := newJSError(err)
					reject.Invoke(jsErr)
					return
				}

				resolve.Invoke(val)
			}()
			return js.Undefined()
		})
		return js.Global().Get("Promise").New(promiseExecutor)
	})
}

func asyncFn(f func(js.Value, []js.Value) (any, error)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		var promiseExecutor js.Func
		promiseExecutor = js.FuncOf(func(_ js.Value, promiseArgs []js.Value) any {
			go func() {
				defer promiseExecutor.Release()
				resolve := promiseArgs[0]
				reject := promiseArgs[1]

				result, err := f(this, args)
				if err != nil {
					jsErr := newJSError(err)
					reject.Invoke(jsErr)
					return
				}
				resolve.Invoke(result)
			}()
			return js.Undefined()
		})
		return js.Global().Get("Promise").New(promiseExecutor)
	})
}
