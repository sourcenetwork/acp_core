// Code generated from ./pkg/parser/theorem_parser/Theorem.g4 by ANTLR 4.13.2. DO NOT EDIT.

package theorem_parser // Theorem
import "github.com/antlr4-go/antlr/v4"


// A complete Visitor for a parse tree produced by TheoremParser.
type TheoremVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by TheoremParser#relationship_document.
	VisitRelationship_document(ctx *Relationship_documentContext) interface{}

	// Visit a parse tree produced by TheoremParser#relationship_set.
	VisitRelationship_set(ctx *Relationship_setContext) interface{}

	// Visit a parse tree produced by TheoremParser#policy_thorem.
	VisitPolicy_thorem(ctx *Policy_thoremContext) interface{}

	// Visit a parse tree produced by TheoremParser#authorization_theorems.
	VisitAuthorization_theorems(ctx *Authorization_theoremsContext) interface{}

	// Visit a parse tree produced by TheoremParser#authorization_theorem.
	VisitAuthorization_theorem(ctx *Authorization_theoremContext) interface{}

	// Visit a parse tree produced by TheoremParser#implied_relations.
	VisitImplied_relations(ctx *Implied_relationsContext) interface{}

	// Visit a parse tree produced by TheoremParser#implied_relation.
	VisitImplied_relation(ctx *Implied_relationContext) interface{}

	// Visit a parse tree produced by TheoremParser#object_rel.
	VisitObject_rel(ctx *Object_relContext) interface{}

	// Visit a parse tree produced by TheoremParser#delegation_theorems.
	VisitDelegation_theorems(ctx *Delegation_theoremsContext) interface{}

	// Visit a parse tree produced by TheoremParser#delegation_theorem.
	VisitDelegation_theorem(ctx *Delegation_theoremContext) interface{}

	// Visit a parse tree produced by TheoremParser#relationship.
	VisitRelationship(ctx *RelationshipContext) interface{}

	// Visit a parse tree produced by TheoremParser#operation.
	VisitOperation(ctx *OperationContext) interface{}

	// Visit a parse tree produced by TheoremParser#subj_uset.
	VisitSubj_uset(ctx *Subj_usetContext) interface{}

	// Visit a parse tree produced by TheoremParser#subj_obj.
	VisitSubj_obj(ctx *Subj_objContext) interface{}

	// Visit a parse tree produced by TheoremParser#subj_actor.
	VisitSubj_actor(ctx *Subj_actorContext) interface{}

	// Visit a parse tree produced by TheoremParser#object.
	VisitObject(ctx *ObjectContext) interface{}

	// Visit a parse tree produced by TheoremParser#ascii_id.
	VisitAscii_id(ctx *Ascii_idContext) interface{}

	// Visit a parse tree produced by TheoremParser#utf_id.
	VisitUtf_id(ctx *Utf_idContext) interface{}

	// Visit a parse tree produced by TheoremParser#relation.
	VisitRelation(ctx *RelationContext) interface{}

	// Visit a parse tree produced by TheoremParser#resource.
	VisitResource(ctx *ResourceContext) interface{}

	// Visit a parse tree produced by TheoremParser#actorid.
	VisitActorid(ctx *ActoridContext) interface{}

}