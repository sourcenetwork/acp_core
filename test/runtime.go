package test

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/stretchr/testify/require"
)

// NewTestRuntime returns a runtime for executing tests
func NewTestRuntime(t *testing.T) runtime.RuntimeManager {
	manager, err := runtime.NewRuntimeManager(runtime.WithMemKV())
	require.Nil(t, err)
	return manager
}
