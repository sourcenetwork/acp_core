package telemetry

import (
	"context"

	"github.com/google/uuid"
	"github.com/sourcenetwork/acp_core/internal/decorator"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

type requestCtxKey int

var requestCtxValue requestCtxKey = 0

// GetRequestContextData returns the context's RequestData
func GetRequestContextData(ctx context.Context) *RequestData {
	return ctx.Value(requestCtxValue).(*RequestData)
}

// RequestData models request specific processing data
type RequestData struct {
	UUID        string
	SandboxData *types.SandboxData
}

// InitRequestContext initializes the request specific data
func InitRequestContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, requestCtxValue, &RequestData{
		UUID: uuid.NewString(),
	})
}

// RequestDataInitializerDecorator adds an instance of RequestData to the ctx
func RequestDataInitializerDecorator(h decorator.Handler) decorator.Handler {
	return func(ctx context.Context, req any) (any, error) {
		ctx = InitRequestContext(ctx)
		return h(ctx, req)
	}
}
