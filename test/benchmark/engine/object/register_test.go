package object

import (
	"fmt"
	"testing"

	"github.com/google/uuid"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
)

var policyStr string = `name: policy
resources:
- name: doc
  relations:
  - name: owner
    types:
    - actor
spec: none
`

func BenchmarkObjectRegistration(b *testing.B) {
	// range from [2^8, 2^20] objs
	for i := 8; i <= 20; i++ {
		var objCount int = 1 << i
		b.Run(fmt.Sprintf("obj=%v", objCount), func(b *testing.B) {
			objs := make([]*types.Object, 0, objCount)
			for j := 0; j < objCount; j++ {
				id := uuid.New().String()
				obj := types.NewObject("doc", id)
				objs = append(objs, obj)

			}

			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				ctx := test.NewTestCtx(b)
				ctx.SetPrincipal("alice")

				action := test.CreatePolicyAction{
					Policy: policyStr,
				}
				pol := action.Run(ctx)

				a1 := test.RegisterObjectsAction{
					PolicyId: pol.Id,
					Objects:  objs,
				}
				a1.Run(ctx)
			}
		})
	}
}
