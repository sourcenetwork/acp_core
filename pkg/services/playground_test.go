package services

import (
	"context"
	"testing"

	"github.com/sourcenetwork/acp_core/internal/decorator"
	"github.com/sourcenetwork/acp_core/internal/telemetry"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/stretchr/testify/require"
)

var _ telemetry.ErrorPublshingClient = (*publishErrorClientMock)(nil)

type publishErrorClientMock struct {
	Called bool
}

func (m *publishErrorClientMock) PushError(ctx context.Context, data *types.SandboxData, err error) error {
	m.Called = true
	return nil
}

func Test_PlaygroundService_DecoratorRecoversFromPanic(t *testing.T) {
	dec := newPlaygroundDecorator(nil, false)

	handler := func(ctx context.Context, _ any) (any, error) {
		panic("test")
	}

	h := decorator.Decorate(handler, dec)
	_, err := h(context.TODO(), nil)
	require.ErrorIs(t, err, errors.ErrorType_INTERNAL)
	require.Contains(t, err.Error(), "panic")
}

func Test_PlaygroundService_ConfiguredDecorator_PushesErrorToBackend(t *testing.T) {
	data := &types.SandboxData{
		PolicyDefinition: "pol",
		Relationships:    "rels",
		PolicyTheorem:    "thm",
	}
	m := publishErrorClientMock{}
	dec := newPlaygroundDecorator(&m, false)

	handler := func(ctx context.Context, _ any) (any, error) {
		d := telemetry.GetRequestContextData(ctx)
		d.SandboxData = data
		panic("test")
	}
	h := decorator.Decorate(handler, dec)

	_, _ = h(context.TODO(), nil)
	require.True(t, m.Called)
}

func Test_PlaygroundService_DecoratorInjectsRequestContextDataEvenWithPanic(t *testing.T) {
	var ctx context.Context
	_ = ctx //ctx is overwritten by handler

	dec := newPlaygroundDecorator(nil, false)
	handler := func(ctxInt context.Context, _ any) (any, error) {
		ctx = ctxInt
		panic("!")
	}
	h := decorator.Decorate(handler, dec)

	_, _ = h(context.TODO(), nil)

	data := telemetry.GetRequestContextData(ctx)
	require.NotNil(t, data)
}
