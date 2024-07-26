package parser

import (
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

type IndexedObject[T any] struct {
	Obj  T
	Line int
	Col  int
}

func NewIndexedObject[T any](obj T, line, col int) IndexedObject[T] {
	return IndexedObject[T]{
		Obj:  obj,
		Line: line,
		Col:  col,
	}
}

type IndexedPolicyTheorem struct {
	DelegationTheorems    []IndexedObject[*types.DelegationTheorem]
	AuthorizationTheorems []IndexedObject[*types.AuthorizationTheorem]
}

func (t *IndexedPolicyTheorem) ToPolicyTheorem() *types.PolicyTheorem {
	return &types.PolicyTheorem{
		DelegationTheorems:    utils.MapSlice(t.DelegationTheorems, func(o IndexedObject[*types.DelegationTheorem]) *types.DelegationTheorem { return o.Obj }),
		AuthorizationTheorems: utils.MapSlice(t.AuthorizationTheorems, func(o IndexedObject[*types.AuthorizationTheorem]) *types.AuthorizationTheorem { return o.Obj }),
	}
}
