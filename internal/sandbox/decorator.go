package sandbox

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/decorator"
	"github.com/sourcenetwork/acp_core/internal/telemetry"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// InjectSandboxData injects data into the context's RequestContext, if it was preivously initialized
func InjectSandboxData(ctx context.Context, data *types.SandboxData) {
	d := decorator.GetRequestContextData(ctx)
	if d != nil {
		d.SandboxData = data
	}
}

// InternalErrorPublisherDecorator publishes internal errors to the playground backend service
func InternalErrorPublisherDecorator(client *telemetry.PlaygroundBackendErrorClient) decorator.Decorator {
	return func(h decorator.Handler) decorator.Handler {
		return func(ctx context.Context, req any) (any, error) {
			resp, err := h(ctx, req)
			if err != nil && errors.Is(err, errors.ErrorType_INTERNAL) {
				data := decorator.GetRequestContextData(ctx)
				if data != nil && data.SandboxData != nil {
					httpErr := client.PushError(ctx, data.SandboxData, err)
					if httpErr != nil {
						// TODO log error
					}
				}
			}
			return resp, err
		}
	}
}
