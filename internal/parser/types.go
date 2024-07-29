package parser

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

type IndexedObject[T any] struct {
	Obj   T
	Range types.BufferRange
}

func NewIndexedObject[T any](obj T, ctx antlr.ParserRuleContext) IndexedObject[T] {
	return IndexedObject[T]{
		Obj: obj,
		Range: types.BufferRange{
			Start: &types.BufferPosition{
				Line:   uint64(ctx.GetStart().GetLine()),
				Column: uint64(ctx.GetStart().GetColumn()),
			},
			End: &types.BufferPosition{
				Line:   uint64(ctx.GetStop().GetLine()),
				Column: uint64(ctx.GetStop().GetColumn()),
			},
		},
	}
}

type IndexedPolicyTheorem struct {
	DelegationTheorems    []IndexedObject[*types.DelegationTheorem]
	AuthorizationTheorems []IndexedObject[*types.AuthorizationTheorem]
	ReachabilityTheorems  []IndexedObject[*types.ReachabilityTheorem]
}

func (t *IndexedPolicyTheorem) ToPolicyTheorem() *types.PolicyTheorem {
	return &types.PolicyTheorem{
		DelegationTheorems:    utils.MapSlice(t.DelegationTheorems, func(o IndexedObject[*types.DelegationTheorem]) *types.DelegationTheorem { return o.Obj }),
		AuthorizationTheorems: utils.MapSlice(t.AuthorizationTheorems, func(o IndexedObject[*types.AuthorizationTheorem]) *types.AuthorizationTheorem { return o.Obj }),
	}
}
