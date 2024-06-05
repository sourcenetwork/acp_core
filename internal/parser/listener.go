package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ TestSuiteListener = &relationshipListener{}

// relationshipListener is a complete listener for a parse tree produced by TestSuiteParser.
type relationshipListener struct {
	relationship *types.Relationship
	objects      []*types.Object
	relations    []string
}

func (l *relationshipListener) GetRelationship() *types.Relationship {
	return l.relationship
}

func (s *relationshipListener) EnterRelationship(ctx *RelationshipContext) {
	s.relationship = &types.Relationship{
		Object:  &types.Object{},
		Subject: &types.Subject{},
	}
}

func (l *relationshipListener) EnterResource(c *ResourceContext) {
	l.objects[len(l.objects)-1].Resource = c.GetText()
}

func (l *relationshipListener) EnterObject(c *ObjectContext) {
	l.objects = append(l.objects, &types.Object{})
}

func (l *relationshipListener) EnterAscii_id(c *Ascii_idContext) {
	l.objects[len(l.objects)-1].Id = c.GetText()
}

func (l *relationshipListener) ExitRelationship(c *RelationshipContext) {
	l.relationship.Object = l.objects[0]
	l.relationship.Relation = l.relations[0]
	l.objects = nil
	l.relations = nil
}

func (l *relationshipListener) EnterUtf_id(c *Utf_idContext) {
	txt := c.GetText()
	id := c.GetText()[1 : len(txt)-1] // remove quotes around string literal
	l.objects[len(l.objects)-1].Id = id
}

func (l *relationshipListener) EnterRelation(c *RelationContext) {
	l.relations = append(l.relations, c.GetText())
}

func (l *relationshipListener) ExitSubj_obj(c *Subj_objContext) {
	len := len(l.objects)
	obj := l.objects[len-1]
	l.objects = l.objects[0 : len-1]
	l.relationship.Subject.Subject = &types.Subject_Object{
		Object: obj,
	}
}

func (l *relationshipListener) ExitSubj_uset(c *Subj_usetContext) {
	objsLen := len(l.objects)
	obj := l.objects[objsLen-1]
	l.objects = l.objects[0 : objsLen-1]

	relsLen := len(l.relations)
	rel := l.relations[relsLen-1]
	l.relations = l.relations[0 : relsLen-1]

	l.relationship.Subject.Subject = &types.Subject_ActorSet{
		ActorSet: &types.ActorSet{
			Object:   obj,
			Relation: rel,
		},
	}
}

func (l *relationshipListener) ExitSubj_actor(c *Subj_actorContext) {
	l.relationship.Subject.Subject = &types.Subject_Actor{
		Actor: &types.Actor{
			Id: c.GetText(),
		},
	}
}

func (l *relationshipListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
}

// VisitTerminal is called when a terminal node is visited.
func (s *relationshipListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *relationshipListener) VisitErrorNode(node antlr.ErrorNode) {}

// ExitEveryRule is called when any rule is exited.
func (s *relationshipListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterRelationship_set is called when production relationship_set is entered.
func (s *relationshipListener) EnterRelationship_set(ctx *Relationship_setContext) {}

// ExitRelationship_set is called when production relationship_set is exited.
func (s *relationshipListener) ExitRelationship_set(ctx *Relationship_setContext) {}

// EnterPolicy_tests is called when production policy_tests is entered.
func (s *relationshipListener) EnterPolicy_tests(ctx *Policy_testsContext) {}

// ExitPolicy_tests is called when production policy_tests is exited.
func (s *relationshipListener) ExitPolicy_tests(ctx *Policy_testsContext) {}

// EnterChecks is called when production checks is entered.
func (s *relationshipListener) EnterChecks(ctx *ChecksContext) {}

// ExitChecks is called when production checks is exited.
func (s *relationshipListener) ExitChecks(ctx *ChecksContext) {}

// EnterCheck is called when production check is entered.
func (s *relationshipListener) EnterCheck(ctx *CheckContext) {}

// ExitCheck is called when production check is exited.
func (s *relationshipListener) ExitCheck(ctx *CheckContext) {}

// EnterImplied_relations is called when production implied_relations is entered.
func (s *relationshipListener) EnterImplied_relations(ctx *Implied_relationsContext) {}

// ExitImplied_relations is called when production implied_relations is exited.
func (s *relationshipListener) ExitImplied_relations(ctx *Implied_relationsContext) {}

// EnterImplied_relation is called when production implied_relation is entered.
func (s *relationshipListener) EnterImplied_relation(ctx *Implied_relationContext) {}

// ExitImplied_relation is called when production implied_relation is exited.
func (s *relationshipListener) ExitImplied_relation(ctx *Implied_relationContext) {}

// EnterObject_rel is called when production object_rel is entered.
func (s *relationshipListener) EnterObject_rel(ctx *Object_relContext) {}

// ExitObject_rel is called when production object_rel is exited.
func (s *relationshipListener) ExitObject_rel(ctx *Object_relContext) {}

// EnterDelegation_assertions is called when production delegation_assertions is entered.
func (s *relationshipListener) EnterDelegation_assertions(ctx *Delegation_assertionsContext) {}

// ExitDelegation_assertions is called when production delegation_assertions is exited.
func (s *relationshipListener) ExitDelegation_assertions(ctx *Delegation_assertionsContext) {}

// EnterDelegation_assertion is called when production delegation_assertion is entered.
func (s *relationshipListener) EnterDelegation_assertion(ctx *Delegation_assertionContext) {}

// ExitDelegation_assertion is called when production delegation_assertion is exited.
func (s *relationshipListener) ExitDelegation_assertion(ctx *Delegation_assertionContext) {}

// EnterSubj_uset is called when production subj_uset is entered.
func (s *relationshipListener) EnterSubj_uset(ctx *Subj_usetContext) {}

// EnterSubj_obj is called when production subj_obj is entered.
func (s *relationshipListener) EnterSubj_obj(ctx *Subj_objContext) {}

// EnterSubj_actor is called when production subj_actor is entered.
func (s *relationshipListener) EnterSubj_actor(ctx *Subj_actorContext) {}

// ExitObject is called when production object is exited.
func (s *relationshipListener) ExitObject(ctx *ObjectContext) {}

// ExitAscii_id is called when production ascii_id is exited.
func (s *relationshipListener) ExitAscii_id(ctx *Ascii_idContext) {}

// ExitUtf_id is called when production utf_id is exited.
func (s *relationshipListener) ExitUtf_id(ctx *Utf_idContext) {}

// ExitRelation is called when production relation is exited.
func (s *relationshipListener) ExitRelation(ctx *RelationContext) {}

// ExitResource is called when production resource is exited.
func (s *relationshipListener) ExitResource(ctx *ResourceContext) {}

// EnterActorid is called when production actorid is entered.
func (s *relationshipListener) EnterActorid(ctx *ActoridContext) {}

// ExitActorid is called when production actorid is exited.
func (s *relationshipListener) ExitActorid(ctx *ActoridContext) {}
