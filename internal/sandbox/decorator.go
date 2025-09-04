package sandbox

import (
	"context"
	"time"

	"github.com/sourcenetwork/acp_core/internal/decorator"
	"github.com/sourcenetwork/acp_core/internal/telemetry"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

const maxPublishRetry = 3

var retryTime = time.Millisecond * 500

// InjectSandboxData injects data into the context's RequestContext, if it was preivously initialized
func InjectSandboxData(ctx context.Context, data *types.SandboxData) {
	d := telemetry.GetRequestContextData(ctx)
	if d != nil {
		d.SandboxData = data
	}
}

// InternalErrorPublisherDecorator publishes internal errors to the playground backend service
//
// Setting async to true will dispatch the create publshing of the error in a new goroutine.
// This should be the default, specially for WASM builds as it can deadlock the wasm module
// See: https://pkg.go.dev/syscall/js#FuncOf
func InternalErrorPublisherDecorator(client telemetry.ErrorPublshingClient, async bool) decorator.Decorator {
	return func(h decorator.Handler) decorator.Handler {
		return func(ctx context.Context, req any) (any, error) {
			resp, err := h(ctx, req)
			if err != nil && errors.Is(err, errors.ErrorType_INTERNAL) {
				state := &types.SandboxData{}
				if data := telemetry.GetRequestContextData(ctx); data != nil && data.SandboxData != nil {
					state = data.SandboxData
				}
				worker := func() {
					var httpErr error
					for i := 0; i < maxPublishRetry; i++ {
						httpErr := client.PushError(context.Background(), state, err)
						if httpErr == nil {
							break
						}
						// exponential backoff
						time.Sleep(retryTime * time.Duration(i))
					}
					if httpErr != nil {
						// TODO log error
					}
				}
				if async {
					go worker()
				} else {
					worker()
				}
			}
			return resp, err
		}
	}
}
