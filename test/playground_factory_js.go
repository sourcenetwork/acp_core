//go:build js

package test

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test/js"
)

// playgroundFactory js returns a PlaygroundJS implementation of the PlagroundService for JS tests
func playgroundFactory(t *testing.T, manager runtime.RuntimeManager) types.PlaygroundServiceServer {
	t.Log("using JS Playground")
	return js.NewPlaygroundJS(t, manager)
}
