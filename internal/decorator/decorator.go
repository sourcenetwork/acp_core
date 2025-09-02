package decorator

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/sourcenetwork/acp_core/pkg/errors"
)

// RecoverDecorator recovers from a panic by catching it and returning an internal error
func RecoverDecorator(h Handler) Handler {
	return func(ctx context.Context, req any) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				msg := fmt.Sprintf("recovered from panic: %v", r)
				stack := debug.Stack()
				err = errors.New(msg, errors.ErrorType_INTERNAL, errors.Pair("stack_trace", string(stack)))
				resp = nil
			}
		}()
		return h(ctx, req)
	}
}

// RequestDataInitializerDecorator adds an instance of RequestData to the ctx
func RequestDataInitializerDecorator(h Handler) Handler {
	return func(ctx context.Context, req any) (any, error) {
		ctx = InitRequestContext(ctx)
		return h(ctx, req)
	}
}
