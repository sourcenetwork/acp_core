package types

import (
	"strings"

	"github.com/cosmos/gogoproto/jsonpb"
)

const ttuOperatorLexeme string = "->"

func (t *PermissionFetchTree) MarshalJSON() (string, error) {
	marshaler := jsonpb.Marshaler{}
	return marshaler.MarshalToString(t)
}

func (t *PermissionFetchTree) UnmarshalJSON(json string) error {
	return jsonpb.UnmarshalString(json, t)
}

// IntoPermissionExpr serializes a permission fetch tree back into a string version of a permission expression
func (t *PermissionFetchTree) IntoPermissionExpr() string {
	return t.intoPermissionExpr(&strings.Builder{})
}

// IntoPermissionExpr serializes a permission fetch tree back into a string version of a permission expression
// uses a string builder as accumulator
func (t *PermissionFetchTree) intoPermissionExpr(builder *strings.Builder) string {
	if t.GetOperation() != nil {
		t.GetOperation().intoPermissionExpr(builder)
	} else {
		combNode := t.GetCombNode()
		builder.WriteRune('(')
		combNode.Left.intoPermissionExpr(builder)
		builder.WriteRune(' ')
		combNode.Combinator.intoPermissionExpr(builder)
		builder.WriteRune(' ')
		combNode.Right.intoPermissionExpr(builder)
		builder.WriteRune(')')
	}
	return builder.String()
}

// IntoPermissionExpr maps a FetchOperation back into the permission expression string form
func (op *FetchOperation) intoPermissionExpr(builder *strings.Builder) {
	if op.GetCu() != nil {
		builder.WriteString(op.GetCu().Relation)
	} else if op.GetTtu() != nil {
		ttu := op.GetTtu()
		builder.WriteString(ttu.Resource)
		builder.WriteString(ttuOperatorLexeme)
		builder.WriteString(ttu.Relation)
	} else {
		builder.WriteString("_this")
	}
}

func (c Combinator) intoPermissionExpr(builder *strings.Builder) {
	switch c {
	case Combinator_DIFFERENCE:
		builder.WriteRune('-')
	case Combinator_INTERSECTION:
		builder.WriteRune('&')
	case Combinator_UNION:
		builder.WriteRune('+')
	}
}
