// Code generated from ./pkg/parser/permission_parser/PermissionExpr.g4 by ANTLR 4.13.2. DO NOT EDIT.

package permission_parser // PermissionExpr
import "github.com/antlr4-go/antlr/v4"


// A complete Visitor for a parse tree produced by PermissionExprParser.
type PermissionExprVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by PermissionExprParser#atom.
	VisitAtom(ctx *AtomContext) interface{}

	// Visit a parse tree produced by PermissionExprParser#nested.
	VisitNested(ctx *NestedContext) interface{}

	// Visit a parse tree produced by PermissionExprParser#cu_term.
	VisitCu_term(ctx *Cu_termContext) interface{}

	// Visit a parse tree produced by PermissionExprParser#ttu_term.
	VisitTtu_term(ctx *Ttu_termContext) interface{}

	// Visit a parse tree produced by PermissionExprParser#expr_term.
	VisitExpr_term(ctx *Expr_termContext) interface{}

	// Visit a parse tree produced by PermissionExprParser#relation.
	VisitRelation(ctx *RelationContext) interface{}

	// Visit a parse tree produced by PermissionExprParser#resource.
	VisitResource(ctx *ResourceContext) interface{}

	// Visit a parse tree produced by PermissionExprParser#union.
	VisitUnion(ctx *UnionContext) interface{}

	// Visit a parse tree produced by PermissionExprParser#difference.
	VisitDifference(ctx *DifferenceContext) interface{}

	// Visit a parse tree produced by PermissionExprParser#intersection.
	VisitIntersection(ctx *IntersectionContext) interface{}

}