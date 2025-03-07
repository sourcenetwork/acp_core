package parser

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

var _ error = (*ParserReport)(nil)
var _ errors.TypedError = (*ParserReport)(nil)

// LocatedObject models a parsed Object and a range
// pointing to the span in the input stream from which the object was parsed
type LocatedObject[T any] struct {
	Obj      T
	Interval *types.BufferInterval
}

// NewLocatedObjectFromCtx creates a new ObjectWithInterval from an ANTLR Parser Context
func NewLocatedObjectFromCtx[T any](obj T, ctx antlr.ParserRuleContext) LocatedObject[T] {
	return LocatedObject[T]{
		Obj: obj,
		Interval: &types.BufferInterval{
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

// ParserReport models a set of messages emitted by the parser
type ParserReport struct {
	msg  string
	msgs []*types.LocatedMessage
}

func NewParserReport(msg string, msgs ...*types.LocatedMessage) *ParserReport {
	return &ParserReport{
		msg:  msg,
		msgs: msgs,
	}
}

func (r *ParserReport) AddLocatedMessage(msg *types.LocatedMessage) {
	r.msgs = append(r.msgs, msg)
}

func (r *ParserReport) GetMessages() []*types.LocatedMessage {
	return r.msgs
}

func (r *ParserReport) HasError() bool {
	errs := utils.FilterSlice(r.msgs, func(msg *types.LocatedMessage) bool { return msg.IsError() })
	return len(errs) != 0
}

// ToMultiError collapses the located errors into flat messsages and collects the error into a multi error
func (r *ParserReport) ToMultiError(baseErr *errors.Error) *errors.MultiError {
	return errors.NewMultiError(baseErr, r.getErrors()...)
}

func (r *ParserReport) getErrors() []error {
	return utils.MapFilterSlice(r.msgs,
		func(msg *types.LocatedMessage) bool { return msg.IsError() },
		func(msg *types.LocatedMessage) error { return msg.ToError() },
	)
}

func (r *ParserReport) Error() string {
	if !r.HasError() {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString(r.msg)
	builder.WriteString("\n")
	for _, err := range r.getErrors() {
		builder.WriteString(err.Error())
		builder.WriteString("\n")
	}
	return builder.String()
}

func (r *ParserReport) GetType() errors.ErrorType {
	return errors.ErrorType_BAD_INPUT
}

func (r *ParserReport) Sucess() bool {
	return !r.HasError()
}
