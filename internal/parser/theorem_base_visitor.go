// Code generated from ./internal/parser/Theorem.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // Theorem
import "github.com/antlr4-go/antlr/v4"

type BaseTheoremVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseTheoremVisitor) VisitRelationship_set(ctx *Relationship_setContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitPolicy_thorem(ctx *Policy_thoremContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitAuthorization_theorems(ctx *Authorization_theoremsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitAuthorization_theorem(ctx *Authorization_theoremContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitImplied_relations(ctx *Implied_relationsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitImplied_relation(ctx *Implied_relationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitObject_rel(ctx *Object_relContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitDelegation_theorems(ctx *Delegation_theoremsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitDelegation_theorem(ctx *Delegation_theoremContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitRelationship(ctx *RelationshipContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitOperation(ctx *OperationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitSubj_uset(ctx *Subj_usetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitSubj_obj(ctx *Subj_objContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitSubj_actor(ctx *Subj_actorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitObject(ctx *ObjectContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitAscii_id(ctx *Ascii_idContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitUtf_id(ctx *Utf_idContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitRelation(ctx *RelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitResource(ctx *ResourceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTheoremVisitor) VisitActorid(ctx *ActoridContext) interface{} {
	return v.VisitChildren(ctx)
}
