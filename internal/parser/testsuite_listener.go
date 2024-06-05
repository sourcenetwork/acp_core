// Code generated from TestSuite.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // TestSuite

import "github.com/antlr4-go/antlr/v4"


// TestSuiteListener is a complete listener for a parse tree produced by TestSuiteParser.
type TestSuiteListener interface {
	antlr.ParseTreeListener

	// EnterRelationship_set is called when entering the relationship_set production.
	EnterRelationship_set(c *Relationship_setContext)

	// EnterPolicy_tests is called when entering the policy_tests production.
	EnterPolicy_tests(c *Policy_testsContext)

	// EnterChecks is called when entering the checks production.
	EnterChecks(c *ChecksContext)

	// EnterCheck is called when entering the check production.
	EnterCheck(c *CheckContext)

	// EnterImplied_relations is called when entering the implied_relations production.
	EnterImplied_relations(c *Implied_relationsContext)

	// EnterImplied_relation is called when entering the implied_relation production.
	EnterImplied_relation(c *Implied_relationContext)

	// EnterObject_rel is called when entering the object_rel production.
	EnterObject_rel(c *Object_relContext)

	// EnterDelegation_assertions is called when entering the delegation_assertions production.
	EnterDelegation_assertions(c *Delegation_assertionsContext)

	// EnterDelegation_assertion is called when entering the delegation_assertion production.
	EnterDelegation_assertion(c *Delegation_assertionContext)

	// EnterRelationship is called when entering the relationship production.
	EnterRelationship(c *RelationshipContext)

	// EnterSubj_uset is called when entering the subj_uset production.
	EnterSubj_uset(c *Subj_usetContext)

	// EnterSubj_obj is called when entering the subj_obj production.
	EnterSubj_obj(c *Subj_objContext)

	// EnterSubj_actor is called when entering the subj_actor production.
	EnterSubj_actor(c *Subj_actorContext)

	// EnterObject is called when entering the object production.
	EnterObject(c *ObjectContext)

	// EnterAscii_id is called when entering the ascii_id production.
	EnterAscii_id(c *Ascii_idContext)

	// EnterUtf_id is called when entering the utf_id production.
	EnterUtf_id(c *Utf_idContext)

	// EnterRelation is called when entering the relation production.
	EnterRelation(c *RelationContext)

	// EnterResource is called when entering the resource production.
	EnterResource(c *ResourceContext)

	// EnterActorid is called when entering the actorid production.
	EnterActorid(c *ActoridContext)

	// ExitRelationship_set is called when exiting the relationship_set production.
	ExitRelationship_set(c *Relationship_setContext)

	// ExitPolicy_tests is called when exiting the policy_tests production.
	ExitPolicy_tests(c *Policy_testsContext)

	// ExitChecks is called when exiting the checks production.
	ExitChecks(c *ChecksContext)

	// ExitCheck is called when exiting the check production.
	ExitCheck(c *CheckContext)

	// ExitImplied_relations is called when exiting the implied_relations production.
	ExitImplied_relations(c *Implied_relationsContext)

	// ExitImplied_relation is called when exiting the implied_relation production.
	ExitImplied_relation(c *Implied_relationContext)

	// ExitObject_rel is called when exiting the object_rel production.
	ExitObject_rel(c *Object_relContext)

	// ExitDelegation_assertions is called when exiting the delegation_assertions production.
	ExitDelegation_assertions(c *Delegation_assertionsContext)

	// ExitDelegation_assertion is called when exiting the delegation_assertion production.
	ExitDelegation_assertion(c *Delegation_assertionContext)

	// ExitRelationship is called when exiting the relationship production.
	ExitRelationship(c *RelationshipContext)

	// ExitSubj_uset is called when exiting the subj_uset production.
	ExitSubj_uset(c *Subj_usetContext)

	// ExitSubj_obj is called when exiting the subj_obj production.
	ExitSubj_obj(c *Subj_objContext)

	// ExitSubj_actor is called when exiting the subj_actor production.
	ExitSubj_actor(c *Subj_actorContext)

	// ExitObject is called when exiting the object production.
	ExitObject(c *ObjectContext)

	// ExitAscii_id is called when exiting the ascii_id production.
	ExitAscii_id(c *Ascii_idContext)

	// ExitUtf_id is called when exiting the utf_id production.
	ExitUtf_id(c *Utf_idContext)

	// ExitRelation is called when exiting the relation production.
	ExitRelation(c *RelationContext)

	// ExitResource is called when exiting the resource production.
	ExitResource(c *ResourceContext)

	// ExitActorid is called when exiting the actorid production.
	ExitActorid(c *ActoridContext)
}
