package telemetry

import (
	"context"
	"testing"

	"github.com/sourcenetwork/acp_core/internal/decorator"
	"github.com/stretchr/testify/require"
)

func Test_InitRequestContext_ReturnsContextWithRequestData(t *testing.T) {
	ctx := context.TODO()
	ctx = InitRequestContext(ctx)

	data := GetRequestContextData(ctx)
	require.NotNil(t, data)
	require.NotEmpty(t, data.UUID)
}

func Test_RequestDataInitializerDecorator_HandlerContextIncludesRequestData(t *testing.T) {
	h := func(ctx context.Context, _ any) (any, error) {
		return ctx, nil
	}
	h = decorator.Decorate(h, RequestDataInitializerDecorator)

	resp, _ := h(context.TODO(), nil)

	ctx := resp.(context.Context)
	data := GetRequestContextData(ctx)
	require.NotNil(t, data)
	require.NotEmpty(t, data.UUID)
}
