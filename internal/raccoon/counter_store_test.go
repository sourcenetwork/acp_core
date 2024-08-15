package raccoon

import (
	"context"
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/stretchr/testify/require"
)

func TestCounterStore(t *testing.T) {
	manager, err := runtime.NewRuntimeManager()
	require.NoError(t, err)

	counter := NewCounterStoreFromRunetimeManager(manager, "")

	ctx := context.TODO()

	i, err := counter.GetNext(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), i)

	err = counter.Increment(ctx)
	require.NoError(t, err)

	i, err = counter.GetNext(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(2), i)
}
