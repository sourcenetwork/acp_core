package decorator

import (
	"context"
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_RecoverDecorator_ReturnsInternalError(t *testing.T) {
	handler := func(ctx context.Context, _ any) (any, error) {
		panic("!")
	}

	handler = RecoverDecorator(handler)
	resp, err := handler(context.TODO(), nil)

	require.Nil(t, resp)
	require.ErrorIs(t, err, errors.ErrorType_INTERNAL)
}
