//go:build !js

package test

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/services"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// playgroundFactory returns the default PlaygroundService implementation for non-js tests
func playgroundFactory(t testing.TB, manager runtime.RuntimeManager) types.PlaygroundServiceServer {
	return services.NewPlaygroundService(manager, nil)
}
