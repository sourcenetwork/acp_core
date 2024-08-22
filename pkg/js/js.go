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
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		playground := newPlaygroundServiceProxy(ctx)
		return playground.asValue()
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
func newPlaygroundServiceProxy(ctx context.Context) *PlaygroundServiceProxy {
	manager, err := runtime.NewRuntimeManager()
	if err != nil {
		panic(err)
	}
	service := services.NewPlaygroundService(manager)

	proxy := &PlaygroundServiceProxy{
		ctx:     ctx,
		manager: manager,
		service: service,
	}
	proxyMap := map[string]js.Func{
		"newSandbox":        proxy.newSandbox(),
		"listSandboxes":     proxy.listSandboxes(),
		"setState":          proxy.setState(),
		"restoreScratchpad": proxy.restoreScratchpad(),
		"getCatalogue":      proxy.getCatalogue(),
		"verifyTheorems":    proxy.verifyTheorems(),
		"simulate":          proxy.simulate(),
		"close":             proxy.close(),
	}
	proxy.proxyMap = proxyMap
	return proxy
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

func (s *PlaygroundServiceProxy) newSandbox() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return invalidNumberOfArgs()
		}
		reqStr := valueToJson(args[0])

		req := &types.NewSandboxRequest{}
		err := jsonpb.UnmarshalString(reqStr, req)
		if err != nil {
			return errToJson(err)
		}

		resp, err := s.service.NewSandbox(s.ctx, req)
		if err != nil {
			return errToJson(err)
		}
		return toJSObject(resp)
	})
}

func (s *PlaygroundServiceProxy) listSandboxes() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return invalidNumberOfArgs()
		}
		reqStr := valueToJson(args[0])

		req := &types.ListSandboxesRequest{}
		err := jsonpb.UnmarshalString(reqStr, req)
		if err != nil {
			return errToJson(err)
		}

		resp, err := s.service.ListSandboxes(s.ctx, req)
		if err != nil {
			return errToJson(err)
		}
		return toJSObject(resp)
	})
}

func (s *PlaygroundServiceProxy) setState() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return invalidNumberOfArgs()
		}
		reqStr := valueToJson(args[0])

		req := &types.SetStateRequest{}
		err := jsonpb.UnmarshalString(reqStr, req)
		if err != nil {
			return errToJson(err)
		}

		resp, err := s.service.SetState(s.ctx, req)
		if err != nil {
			return errToJson(err)
		}
		return toJSObject(resp)
	})
}

func (s *PlaygroundServiceProxy) restoreScratchpad() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return invalidNumberOfArgs()
		}
		reqStr := valueToJson(args[0])

		req := &types.RestoreScratchpadRequest{}
		err := jsonpb.UnmarshalString(reqStr, req)
		if err != nil {
			return errToJson(err)
		}

		resp, err := s.service.RestoreScratchpad(s.ctx, req)
		if err != nil {
			return errToJson(err)
		}
		return toJSObject(resp)
	})
}

func (s *PlaygroundServiceProxy) getCatalogue() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return invalidNumberOfArgs()
		}
		reqStr := valueToJson(args[0])

		req := &types.GetCatalogueRequest{}
		err := jsonpb.UnmarshalString(reqStr, req)
		if err != nil {
			return errToJson(err)
		}

		resp, err := s.service.GetCatalogue(s.ctx, req)
		if err != nil {
			return errToJson(err)
		}
		return toJSObject(resp)
	})
}

func (s *PlaygroundServiceProxy) verifyTheorems() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return invalidNumberOfArgs()
		}
		reqStr := valueToJson(args[0])

		req := &types.VerifyTheoremsRequest{}
		err := jsonpb.UnmarshalString(reqStr, req)
		if err != nil {
			return errToJson(err)
		}

		resp, err := s.service.VerifyTheorems(s.ctx, req)
		if err != nil {
			return errToJson(err)
		}
		return toJSObject(resp)
	})
}

func (s *PlaygroundServiceProxy) simulate() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return invalidNumberOfArgs()
		}
		reqStr := valueToJson(args[0])

		req := &types.SimulateRequest{}
		err := jsonpb.UnmarshalString(reqStr, req)
		if err != nil {
			return errToJson(err)
		}

		resp, err := s.service.Simulate(s.ctx, req)
		if err != nil {
			return errToJson(err)
		}
		return toJSObject(resp)
	})
}

func (s *PlaygroundServiceProxy) close() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		for _, f := range s.proxyMap {
			f.Release()
		}
		return nil
	})
}

func valueToJson(v js.Value) string {
	result := js.Global().Get("JSON").Call("stringify", v)
	if result.Type() != js.TypeString {
		panic("expected string")
	}

	return result.String()
}

func toJSObject[T proto.Message](val T) js.Value {
	marshaler := jsonpb.Marshaler{}
	valStr, err := marshaler.MarshalToString(val)
	if err != nil {
		panic(err)
	}
	jsVal := js.Global().Get("JSON").Call("parse", valStr)
	return jsVal
}

func errToJson(err error) js.Value {
	errMap := map[string]any{
		"message": err.Error(),
	}
	jsVal := js.ValueOf(errMap)
	return jsVal
}

func invalidNumberOfArgs() js.Value {
	errMap := map[string]any{
		"message": "expected 1 argument",
	}
	jsVal := js.ValueOf(errMap)
	return jsVal
}

func newInvalidArgsErr(count int) error {
	return errors.Wrap("invalid number of arguments", errors.ErrorType_BAD_INPUT, errors.Pair("count", count))
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
