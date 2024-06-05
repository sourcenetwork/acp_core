// Code generated from TestSuite.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // TestSuite

import "github.com/antlr4-go/antlr/v4"

// BaseTestSuiteListener is a complete listener for a parse tree produced by TestSuiteParser.
type BaseTestSuiteListener struct{}

var _ TestSuiteListener = &BaseTestSuiteListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseTestSuiteListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseTestSuiteListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseTestSuiteListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseTestSuiteListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterRelationship_set is called when production relationship_set is entered.
func (s *BaseTestSuiteListener) EnterRelationship_set(ctx *Relationship_setContext) {}

// ExitRelationship_set is called when production relationship_set is exited.
func (s *BaseTestSuiteListener) ExitRelationship_set(ctx *Relationship_setContext) {}

// EnterPolicy_tests is called when production policy_tests is entered.
func (s *BaseTestSuiteListener) EnterPolicy_tests(ctx *Policy_testsContext) {}

// ExitPolicy_tests is called when production policy_tests is exited.
func (s *BaseTestSuiteListener) ExitPolicy_tests(ctx *Policy_testsContext) {}

// EnterChecks is called when production checks is entered.
func (s *BaseTestSuiteListener) EnterChecks(ctx *ChecksContext) {}

// ExitChecks is called when production checks is exited.
func (s *BaseTestSuiteListener) ExitChecks(ctx *ChecksContext) {}

// EnterCheck is called when production check is entered.
func (s *BaseTestSuiteListener) EnterCheck(ctx *CheckContext) {}

// ExitCheck is called when production check is exited.
func (s *BaseTestSuiteListener) ExitCheck(ctx *CheckContext) {}

// EnterImplied_relations is called when production implied_relations is entered.
func (s *BaseTestSuiteListener) EnterImplied_relations(ctx *Implied_relationsContext) {}

// ExitImplied_relations is called when production implied_relations is exited.
func (s *BaseTestSuiteListener) ExitImplied_relations(ctx *Implied_relationsContext) {}

// EnterImplied_relation is called when production implied_relation is entered.
func (s *BaseTestSuiteListener) EnterImplied_relation(ctx *Implied_relationContext) {}

// ExitImplied_relation is called when production implied_relation is exited.
func (s *BaseTestSuiteListener) ExitImplied_relation(ctx *Implied_relationContext) {}

// EnterObject_rel is called when production object_rel is entered.
func (s *BaseTestSuiteListener) EnterObject_rel(ctx *Object_relContext) {}

// ExitObject_rel is called when production object_rel is exited.
func (s *BaseTestSuiteListener) ExitObject_rel(ctx *Object_relContext) {}

// EnterDelegation_assertions is called when production delegation_assertions is entered.
func (s *BaseTestSuiteListener) EnterDelegation_assertions(ctx *Delegation_assertionsContext) {}

// ExitDelegation_assertions is called when production delegation_assertions is exited.
func (s *BaseTestSuiteListener) ExitDelegation_assertions(ctx *Delegation_assertionsContext) {}

// EnterDelegation_assertion is called when production delegation_assertion is entered.
func (s *BaseTestSuiteListener) EnterDelegation_assertion(ctx *Delegation_assertionContext) {}

// ExitDelegation_assertion is called when production delegation_assertion is exited.
func (s *BaseTestSuiteListener) ExitDelegation_assertion(ctx *Delegation_assertionContext) {}

// EnterRelationship is called when production relationship is entered.
func (s *BaseTestSuiteListener) EnterRelationship(ctx *RelationshipContext) {}

// ExitRelationship is called when production relationship is exited.
func (s *BaseTestSuiteListener) ExitRelationship(ctx *RelationshipContext) {}

// EnterSubj_uset is called when production subj_uset is entered.
func (s *BaseTestSuiteListener) EnterSubj_uset(ctx *Subj_usetContext) {}

// ExitSubj_uset is called when production subj_uset is exited.
func (s *BaseTestSuiteListener) ExitSubj_uset(ctx *Subj_usetContext) {}

// EnterSubj_obj is called when production subj_obj is entered.
func (s *BaseTestSuiteListener) EnterSubj_obj(ctx *Subj_objContext) {}

// ExitSubj_obj is called when production subj_obj is exited.
func (s *BaseTestSuiteListener) ExitSubj_obj(ctx *Subj_objContext) {}

// EnterSubj_actor is called when production subj_actor is entered.
func (s *BaseTestSuiteListener) EnterSubj_actor(ctx *Subj_actorContext) {}

// ExitSubj_actor is called when production subj_actor is exited.
func (s *BaseTestSuiteListener) ExitSubj_actor(ctx *Subj_actorContext) {}

// EnterObject is called when production object is entered.
func (s *BaseTestSuiteListener) EnterObject(ctx *ObjectContext) {}

// ExitObject is called when production object is exited.
func (s *BaseTestSuiteListener) ExitObject(ctx *ObjectContext) {}

// EnterAscii_id is called when production ascii_id is entered.
func (s *BaseTestSuiteListener) EnterAscii_id(ctx *Ascii_idContext) {}

// ExitAscii_id is called when production ascii_id is exited.
func (s *BaseTestSuiteListener) ExitAscii_id(ctx *Ascii_idContext) {}

// EnterUtf_id is called when production utf_id is entered.
func (s *BaseTestSuiteListener) EnterUtf_id(ctx *Utf_idContext) {}

// ExitUtf_id is called when production utf_id is exited.
func (s *BaseTestSuiteListener) ExitUtf_id(ctx *Utf_idContext) {}

// EnterRelation is called when production relation is entered.
func (s *BaseTestSuiteListener) EnterRelation(ctx *RelationContext) {}

// ExitRelation is called when production relation is exited.
func (s *BaseTestSuiteListener) ExitRelation(ctx *RelationContext) {}

// EnterResource is called when production resource is entered.
func (s *BaseTestSuiteListener) EnterResource(ctx *ResourceContext) {}

// ExitResource is called when production resource is exited.
func (s *BaseTestSuiteListener) ExitResource(ctx *ResourceContext) {}

// EnterActorid is called when production actorid is entered.
func (s *BaseTestSuiteListener) EnterActorid(ctx *ActoridContext) {}

// ExitActorid is called when production actorid is exited.
func (s *BaseTestSuiteListener) ExitActorid(ctx *ActoridContext) {}
