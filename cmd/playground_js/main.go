//go:build js

// package main defines the entrypoint for the playground js
//
// This package was designed to be compiled as a WASM module which exposes an implementation of PlaygrounService.
// The service is exported to the JS environment by mutating the Global JS object, in a browser context this is equivalent to the `window`.
package main

import (
	"context"
	"syscall/js"

	acpJs "github.com/sourcenetwork/acp_core/pkg/js"
)

func main() {
	ctx := context.Background()

	js.Global().Set("AcpPlayground", map[string]any{
		"new": acpJs.PlaygroundConstructor(ctx),
	})

	<-ctx.Done()
}
