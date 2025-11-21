package relationship

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

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
  - name: reader
    types:
    - actor
spec: none
`

func BenchmarkCreateRelationshipsWithObjects(b *testing.B) {
	// range from power of 8 until 19
	for i := 8; i < 21; i++ {
		var objCount int = 1 << i
		var relCount int = 1 << (i - 1)
		b.Run(fmt.Sprintf("obj=%v/rels=%v", objCount, relCount), func(b *testing.B) {
			createNRelationshipsWithMRegisteredObjectsToNewActors(b, relCount, objCount)
		})
	}
}

func createNRelationshipsWithMRegisteredObjectsToNewActors(b *testing.B, relCount, objCount int) {
	for n := 0; n < b.N; n++ {
		// Given registered objects
		ctx := test.NewTestCtx(b)
		ctx.SetPrincipal("alice")
		action := test.CreatePolicyAction{
			Policy: policyStr,
		}
		pol := action.Run(ctx)
		objs := make([]*types.Object, 0, objCount)
		for i := 0; i < objCount; i++ {
			id := uuid.New().String()
			obj := types.NewObject("doc", id)
			objs = append(objs, obj)
		}
		a1 := test.RegisterObjectsAction{
			PolicyId: pol.Id,
			Objects:  objs,
		}
		a1.Run(ctx)

		// When I register relationships to these objects to a new actor
		for i := 0; i < relCount; i++ {
			actorId := "did:example:" + uuid.New().String()
			idx := int(rand.Float64() * float64(objCount))
			rel := types.NewActorRelationship("doc", objs[idx].Id, "reader", actorId)
			_, err := ctx.Engine.SetRelationship(ctx, &types.SetRelationshipRequest{
				PolicyId:     pol.Id,
				Relationship: rel,
			})
			require.NoError(b, err)
		}
	}
}
