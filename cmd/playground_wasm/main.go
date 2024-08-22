//go:build js

package main

import (
	"context"
	"syscall/js"

	acpJs "github.com/sourcenetwork/acp_core/pkg/js"
)

func main() {
	ctx := context.Background()

	js.Global().Set("AcpPlayground", map[string]any{
		"new": acpJs.NewPlayground(ctx),
	})

	<-ctx.Done()
}
