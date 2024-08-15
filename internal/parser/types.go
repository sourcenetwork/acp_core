package parser

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

// LocatedObject models a parsed Object and a range
// pointing to the span in the input stream from which the object was parsed
type LocatedObject[T any] struct {
	Obj   T
	Range types.BufferRange
}

// NewLocatedObjectFromCtx creates a new ObjectWithRange from an ANTLR Parser Context
func NewLocatedObjectFromCtx[T any](obj T, ctx antlr.ParserRuleContext) LocatedObject[T] {
	return LocatedObject[T]{
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

// LocatedPolicyTheorem stores the elements of a PolicyTheorem with
// relevant data which locates each Theorem back to the input stream.
type LocatedPolicyTheorem struct {
	DelegationTheorems    []LocatedObject[*types.DelegationTheorem]
	AuthorizationTheorems []LocatedObject[*types.AuthorizationTheorem]
	ReachabilityTheorems  []LocatedObject[*types.ReachabilityTheorem]
}

// ToPolicyTheorem removes the Location data and returns a PolicyTheorem
func (t *LocatedPolicyTheorem) ToPolicyTheorem() *types.PolicyTheorem {
	return &types.PolicyTheorem{
		DelegationTheorems:    utils.MapSlice(t.DelegationTheorems, func(o LocatedObject[*types.DelegationTheorem]) *types.DelegationTheorem { return o.Obj }),
		AuthorizationTheorems: utils.MapSlice(t.AuthorizationTheorems, func(o LocatedObject[*types.AuthorizationTheorem]) *types.AuthorizationTheorem { return o.Obj }),
	}
}
