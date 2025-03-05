package theorem_parser

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/sourcenetwork/acp_core/internal/parser"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

var _ TheoremVisitor = (*theoremVisitorImpl)(nil)

// theoremVisitorImpl implements TheoremVisitor and walks the parse tree to produce relationships and policy theorems
type theoremVisitorImpl struct {
	*antlr.BaseParseTreeVisitor
}

func (l *theoremVisitorImpl) VisitRelationship_set(ctx *Relationship_setContext) any {
	return utils.MapSlice(ctx.AllRelationship(), func(ctx IRelationshipContext) parser.LocatedObject[*types.Relationship] {
		return l.Visit(ctx).(parser.LocatedObject[*types.Relationship])
	})
}

func (l *theoremVisitorImpl) VisitAscii_id(c *Ascii_idContext) any {
	return c.GetText()
}

func (l *theoremVisitorImpl) VisitObject(ctx *ObjectContext) any {
	resource := ctx.Resource().GetText()
	objId := l.Visit(ctx.Object_id()).(string)
	return types.NewObject(resource, objId)
}

func (l *theoremVisitorImpl) VisitUtf_id(c *Utf_idContext) any {
	txt := c.GetText()
	id := c.GetText()[1 : len(txt)-1] // remove quotes around string literal
	return id
}

func (l *theoremVisitorImpl) VisitRelationship(ctx *RelationshipContext) any {
	rel := &types.Relationship{
		Object:   l.Visit(ctx.Object()).(*types.Object),
		Relation: ctx.Relation().GetText(),
		Subject:  l.Visit(ctx.Subject()).(*types.Subject),
	}
	return parser.NewLocatedObjectFromCtx(rel, ctx)
}

func (l *theoremVisitorImpl) VisitSubj_obj(ctx *Subj_objContext) any {
	obj := l.Visit(ctx.Object()).(*types.Object)
	return &types.Subject{
		Subject: &types.Subject_Object{
			Object: obj,
		},
	}
}

func (l *theoremVisitorImpl) VisitSubj_uset(c *Subj_usetContext) any {
	obj := l.Visit(c.Object()).(*types.Object)
	rel := c.Relation().GetText()
	return &types.Subject{
		Subject: &types.Subject_ActorSet{
			ActorSet: &types.ActorSet{
				Object:   obj,
				Relation: rel,
			},
		},
	}
}

func (l *theoremVisitorImpl) VisitSubj_actor(c *Subj_actorContext) any {
	return &types.Subject{
		Subject: &types.Subject_Actor{
			Actor: &types.Actor{
				Id: c.GetText(),
			},
		},
	}
}

func (l *theoremVisitorImpl) VisitOperation(ctx *OperationContext) any {
	obj := l.Visit(ctx.Object()).(*types.Object)
	rel := ctx.Relation().GetText()
	return &types.Operation{
		Object:     obj,
		Permission: rel,
	}
}

func (l *theoremVisitorImpl) VisitDelegation_theorem(ctx *Delegation_theoremContext) any {
	actor := &types.Actor{Id: ctx.Actorid().GetText()}
	operation := l.Visit(ctx.Operation()).(*types.Operation)
	negate := ctx.NEGATION() != nil
	theorem := &types.DelegationTheorem{
		Actor:      actor,
		Operation:  operation,
		AssertTrue: !negate,
	}
	return parser.NewLocatedObjectFromCtx(theorem, ctx)
}

func (l *theoremVisitorImpl) VisitDelegation_theorems(ctx *Delegation_theoremsContext) any {
	return utils.MapSlice(ctx.AllDelegation_theorem(), func(ctx IDelegation_theoremContext) parser.LocatedObject[*types.DelegationTheorem] {
		return l.Visit(ctx).(parser.LocatedObject[*types.DelegationTheorem])
	})
}

func (l *theoremVisitorImpl) VisitAuthorization_theorem(ctx *Authorization_theoremContext) any {
	negate := ctx.NEGATION() != nil
	relationship := l.Visit(ctx.Relationship()).(parser.LocatedObject[*types.Relationship]).Obj
	theorem := &types.AuthorizationTheorem{
		Operation: &types.Operation{
			Object:     relationship.Object,
			Permission: relationship.Relation,
		},
		Actor:      relationship.GetSubject().GetActor(),
		AssertTrue: !negate,
	}
	return parser.NewLocatedObjectFromCtx(theorem, ctx)
}

func (l *theoremVisitorImpl) VisitAuthorization_theorems(ctx *Authorization_theoremsContext) any {
	return utils.MapSlice(ctx.AllAuthorization_theorem(), func(ctx IAuthorization_theoremContext) parser.LocatedObject[*types.AuthorizationTheorem] {
		return l.Visit(ctx).(parser.LocatedObject[*types.AuthorizationTheorem])
	})
}

func (l *theoremVisitorImpl) VisitPolicy_thorem(ctx *Policy_thoremContext) any {
	authorizationThms := l.Visit(ctx.Authorization_theorems()).([]parser.LocatedObject[*types.AuthorizationTheorem])
	delegationThms := l.Visit(ctx.Delegation_theorems()).([]parser.LocatedObject[*types.DelegationTheorem])
	return &parser.LocatedPolicyTheorem{
		AuthorizationTheorems: authorizationThms,
		DelegationTheorems:    delegationThms,
	}
}

func (v *theoremVisitorImpl) VisitImplied_relations(ctx *Implied_relationsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *theoremVisitorImpl) VisitImplied_relation(ctx *Implied_relationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *theoremVisitorImpl) VisitObject_rel(ctx *Object_relContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *theoremVisitorImpl) VisitRelation(ctx *RelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *theoremVisitorImpl) VisitResource(ctx *ResourceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *theoremVisitorImpl) VisitActorid(ctx *ActoridContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *theoremVisitorImpl) VisitRelationship_document(ctx *Relationship_documentContext) any {
	result := v.Visit(ctx.Relationship())
	return result.(parser.LocatedObject[*types.Relationship])
}

func (v *theoremVisitorImpl) Visit(tree antlr.ParseTree) interface{}         { return tree.Accept(v) }
func (v *theoremVisitorImpl) VisitChildren(_ antlr.RuleNode) interface{}     { return nil }
func (v *theoremVisitorImpl) VisitTerminal(_ antlr.TerminalNode) interface{} { return nil }
func (v *theoremVisitorImpl) VisitErrorNode(_ antlr.ErrorNode) interface{}   { return nil }
