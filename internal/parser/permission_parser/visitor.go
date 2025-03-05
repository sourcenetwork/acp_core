package permission_parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ PermissionExprVisitor = (*visitor)(nil)

type visitor struct {
	antlr.ParseTreeVisitor
}

// Visit a parse tree produced by PermissionExprParser#ttu_term.
// Returns PermissionFetchTree
func (v *visitor) VisitTtu_term(ctx *Ttu_termContext) any {
	return &types.PermissionFetchTree{
		Term: &types.PermissionFetchTree_Operation{
			Operation: &types.FetchOperation{
				Operation: &types.FetchOperation_Ttu{
					Ttu: &types.TupleToUsersetNode{
						Resource: ctx.Resource().GetText(),
						Relation: ctx.Relation().GetText(),
					},
				},
			},
		},
	}
}

// Visit a parse tree produced by PermissionExprParser#cu_term.
func (v *visitor) VisitCu_term(ctx *Cu_termContext) any {
	return &types.PermissionFetchTree{
		Term: &types.PermissionFetchTree_Operation{
			Operation: &types.FetchOperation{
				Operation: &types.FetchOperation_Cu{
					Cu: &types.ComputedUsersetNode{
						Relation: ctx.Relation().GetText(),
					},
				},
			},
		},
	}
}

// Visit a parse tree produced by PermissionExprParser#expr_term.
func (v *visitor) VisitExpr_term(ctx *Expr_termContext) any {
	return v.Visit(ctx.Expr())
}

// Visit a parse tree produced by PermissionExprParser#relation.
func (v *visitor) VisitRelation(ctx *RelationContext) any {
	return ctx.GetText()
}

// Visit a parse tree produced by PermissionExprParser#resource.
func (v *visitor) VisitResource(ctx *ResourceContext) any {
	return ctx.GetText()
}

// Visit a parse tree produced by PermissionExprParser#union.
func (v *visitor) VisitUnion(ctx *UnionContext) any {
	return types.Combinator_UNION
}

// Visit a parse tree produced by PermissionExprParser#difference.
func (v *visitor) VisitDifference(ctx *DifferenceContext) any {
	return types.Combinator_DIFFERENCE
}

// Visit a parse tree produced by PermissionExprParser#intersection.
func (v *visitor) VisitIntersection(ctx *IntersectionContext) any {
	return types.Combinator_INTERSECTION
}

// Visit a parse tree produced by PermissionExprParser#atom.
// return PermissionFetchTree
func (v *visitor) VisitAtom(ctx *AtomContext) any {
	return v.Visit(ctx.Term())
}

// Visit a parse tree produced by PermissionExprParser#nested.
func (v *visitor) VisitNested(ctx *NestedContext) any {
	combinator := v.Visit(ctx.Operator()).(types.Combinator)
	right := v.Visit(ctx.Term()).(*types.PermissionFetchTree)
	left := v.Visit(ctx.Expr()).(*types.PermissionFetchTree)
	return &types.PermissionFetchTree{
		Term: &types.PermissionFetchTree_CombNode{
			CombNode: &types.CombinationNode{
				Left:       left,
				Combinator: combinator,
				Right:      right,
			},
		},
	}
}

func (v *visitor) Visit(tree antlr.ParseTree) any            { return tree.Accept(v) }
func (v *visitor) VisitChildren(node antlr.RuleNode) any     { return nil }
func (v *visitor) VisitTerminal(node antlr.TerminalNode) any { return nil }
func (v *visitor) VisitErrorNode(node antlr.ErrorNode) any   { return nil }
