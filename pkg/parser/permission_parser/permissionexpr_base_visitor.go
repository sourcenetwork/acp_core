// Code generated from ./internal/parser/permission_parser/PermissionExpr.g4 by ANTLR 4.13.2. DO NOT EDIT.

package permission_parser // PermissionExpr
import "github.com/antlr4-go/antlr/v4"

type BasePermissionExprVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BasePermissionExprVisitor) VisitAtom(ctx *AtomContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePermissionExprVisitor) VisitNested(ctx *NestedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePermissionExprVisitor) VisitCu_term(ctx *Cu_termContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePermissionExprVisitor) VisitTtu_term(ctx *Ttu_termContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePermissionExprVisitor) VisitExpr_term(ctx *Expr_termContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePermissionExprVisitor) VisitRelation(ctx *RelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePermissionExprVisitor) VisitResource(ctx *ResourceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePermissionExprVisitor) VisitUnion(ctx *UnionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePermissionExprVisitor) VisitDifference(ctx *DifferenceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePermissionExprVisitor) VisitIntersection(ctx *IntersectionContext) interface{} {
	return v.VisitChildren(ctx)
}
