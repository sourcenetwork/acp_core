package main

import (
	"context"
	"log"

	"github.com/sourcenetwork/acp_core/pkg/engine"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func main() {
	runtime, err := runtime.NewRuntimeManager()
	if err != nil {
		log.Fatal(err)
	}
	engine := engine.NewACPEngine(runtime)
	pol, _ := engine.CreatePolicy(context.TODO(), &types.CreatePolicyRequest{
		Policy:      "",
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	})
	_ = pol
}
