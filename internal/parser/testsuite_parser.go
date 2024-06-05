// Code generated from TestSuite.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // TestSuite

import (
	"fmt"
	"strconv"
  	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}


type TestSuiteParser struct {
	*antlr.BaseParser
}

var TestSuiteParserStaticData struct {
  once                   sync.Once
  serializedATN          []int32
  LiteralNames           []string
  SymbolicNames          []string
  RuleNames              []string
  PredictionContextCache *antlr.PredictionContextCache
  atn                    *antlr.ATN
  decisionToDFA          []*antlr.DFA
}

func testsuiteParserInit() {
  staticData := &TestSuiteParserStaticData
  staticData.LiteralNames = []string{
    "", "'Checks'", "'{'", "'}'", "'ImpliedRelations'", "'=>'", "'#'", "'DelegationAssertions'", 
    "'@'", "':'", "'!'",
  }
  staticData.SymbolicNames = []string{
    "", "", "", "", "", "", "", "", "", "", "NEGATION", "OPERATION", "ID", 
    "STRING", "DID", "HEX", "HEXDIG", "COMMENT", "WS", "NL",
  }
  staticData.RuleNames = []string{
    "relationship_set", "policy_tests", "checks", "check", "implied_relations", 
    "implied_relation", "object_rel", "delegation_assertions", "delegation_assertion", 
    "relationship", "subject", "object", "object_id", "relation", "resource", 
    "actorid",
  }
  staticData.PredictionContextCache = antlr.NewPredictionContextCache()
  staticData.serializedATN = []int32{
	4, 1, 19, 124, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7, 
	4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7, 
	10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15, 
	1, 0, 4, 0, 34, 8, 0, 11, 0, 12, 0, 35, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 
	2, 1, 2, 5, 2, 45, 8, 2, 10, 2, 12, 2, 48, 9, 2, 1, 2, 1, 2, 1, 3, 1, 3, 
	1, 3, 3, 3, 55, 8, 3, 1, 4, 1, 4, 1, 4, 5, 4, 60, 8, 4, 10, 4, 12, 4, 63, 
	9, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 
	1, 7, 1, 7, 5, 7, 78, 8, 7, 10, 7, 12, 7, 81, 9, 7, 1, 7, 1, 7, 1, 8, 1, 
	8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 3, 8, 94, 8, 8, 1, 9, 1, 9, 
	1, 9, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 3, 10, 
	108, 8, 10, 1, 11, 1, 11, 1, 11, 1, 11, 1, 12, 1, 12, 3, 12, 116, 8, 12, 
	1, 13, 1, 13, 1, 14, 1, 14, 1, 15, 1, 15, 1, 15, 0, 0, 16, 0, 2, 4, 6, 
	8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 0, 0, 116, 0, 33, 1, 0, 
	0, 0, 2, 37, 1, 0, 0, 0, 4, 41, 1, 0, 0, 0, 6, 54, 1, 0, 0, 0, 8, 56, 1, 
	0, 0, 0, 10, 66, 1, 0, 0, 0, 12, 70, 1, 0, 0, 0, 14, 74, 1, 0, 0, 0, 16, 
	93, 1, 0, 0, 0, 18, 95, 1, 0, 0, 0, 20, 107, 1, 0, 0, 0, 22, 109, 1, 0, 
	0, 0, 24, 115, 1, 0, 0, 0, 26, 117, 1, 0, 0, 0, 28, 119, 1, 0, 0, 0, 30, 
	121, 1, 0, 0, 0, 32, 34, 3, 18, 9, 0, 33, 32, 1, 0, 0, 0, 34, 35, 1, 0, 
	0, 0, 35, 33, 1, 0, 0, 0, 35, 36, 1, 0, 0, 0, 36, 1, 1, 0, 0, 0, 37, 38, 
	3, 4, 2, 0, 38, 39, 3, 8, 4, 0, 39, 40, 3, 14, 7, 0, 40, 3, 1, 0, 0, 0, 
	41, 42, 5, 1, 0, 0, 42, 46, 5, 2, 0, 0, 43, 45, 3, 6, 3, 0, 44, 43, 1, 
	0, 0, 0, 45, 48, 1, 0, 0, 0, 46, 44, 1, 0, 0, 0, 46, 47, 1, 0, 0, 0, 47, 
	49, 1, 0, 0, 0, 48, 46, 1, 0, 0, 0, 49, 50, 5, 3, 0, 0, 50, 5, 1, 0, 0, 
	0, 51, 55, 3, 18, 9, 0, 52, 53, 5, 10, 0, 0, 53, 55, 3, 18, 9, 0, 54, 51, 
	1, 0, 0, 0, 54, 52, 1, 0, 0, 0, 55, 7, 1, 0, 0, 0, 56, 57, 5, 4, 0, 0, 
	57, 61, 5, 2, 0, 0, 58, 60, 3, 10, 5, 0, 59, 58, 1, 0, 0, 0, 60, 63, 1, 
	0, 0, 0, 61, 59, 1, 0, 0, 0, 61, 62, 1, 0, 0, 0, 62, 64, 1, 0, 0, 0, 63, 
	61, 1, 0, 0, 0, 64, 65, 5, 3, 0, 0, 65, 9, 1, 0, 0, 0, 66, 67, 3, 12, 6, 
	0, 67, 68, 5, 5, 0, 0, 68, 69, 3, 12, 6, 0, 69, 11, 1, 0, 0, 0, 70, 71, 
	3, 22, 11, 0, 71, 72, 5, 6, 0, 0, 72, 73, 3, 26, 13, 0, 73, 13, 1, 0, 0, 
	0, 74, 75, 5, 7, 0, 0, 75, 79, 5, 2, 0, 0, 76, 78, 3, 16, 8, 0, 77, 76, 
	1, 0, 0, 0, 78, 81, 1, 0, 0, 0, 79, 77, 1, 0, 0, 0, 79, 80, 1, 0, 0, 0, 
	80, 82, 1, 0, 0, 0, 81, 79, 1, 0, 0, 0, 82, 83, 5, 3, 0, 0, 83, 15, 1, 
	0, 0, 0, 84, 85, 3, 30, 15, 0, 85, 86, 5, 11, 0, 0, 86, 87, 3, 18, 9, 0, 
	87, 94, 1, 0, 0, 0, 88, 89, 5, 10, 0, 0, 89, 90, 3, 30, 15, 0, 90, 91, 
	5, 11, 0, 0, 91, 92, 3, 18, 9, 0, 92, 94, 1, 0, 0, 0, 93, 84, 1, 0, 0, 
	0, 93, 88, 1, 0, 0, 0, 94, 17, 1, 0, 0, 0, 95, 96, 3, 22, 11, 0, 96, 97, 
	5, 6, 0, 0, 97, 98, 3, 26, 13, 0, 98, 99, 5, 8, 0, 0, 99, 100, 3, 20, 10, 
	0, 100, 19, 1, 0, 0, 0, 101, 102, 3, 22, 11, 0, 102, 103, 5, 6, 0, 0, 103, 
	104, 3, 26, 13, 0, 104, 108, 1, 0, 0, 0, 105, 108, 3, 22, 11, 0, 106, 108, 
	3, 30, 15, 0, 107, 101, 1, 0, 0, 0, 107, 105, 1, 0, 0, 0, 107, 106, 1, 
	0, 0, 0, 108, 21, 1, 0, 0, 0, 109, 110, 3, 28, 14, 0, 110, 111, 5, 9, 0, 
	0, 111, 112, 3, 24, 12, 0, 112, 23, 1, 0, 0, 0, 113, 116, 5, 12, 0, 0, 
	114, 116, 5, 13, 0, 0, 115, 113, 1, 0, 0, 0, 115, 114, 1, 0, 0, 0, 116, 
	25, 1, 0, 0, 0, 117, 118, 5, 12, 0, 0, 118, 27, 1, 0, 0, 0, 119, 120, 5, 
	12, 0, 0, 120, 29, 1, 0, 0, 0, 121, 122, 5, 14, 0, 0, 122, 31, 1, 0, 0, 
	0, 8, 35, 46, 54, 61, 79, 93, 107, 115,
}
  deserializer := antlr.NewATNDeserializer(nil)
  staticData.atn = deserializer.Deserialize(staticData.serializedATN)
  atn := staticData.atn
  staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
  decisionToDFA := staticData.decisionToDFA
  for index, state := range atn.DecisionToState {
    decisionToDFA[index] = antlr.NewDFA(state, index)
  }
}

// TestSuiteParserInit initializes any static state used to implement TestSuiteParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewTestSuiteParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func TestSuiteParserInit() {
  staticData := &TestSuiteParserStaticData
  staticData.once.Do(testsuiteParserInit)
}

// NewTestSuiteParser produces a new parser instance for the optional input antlr.TokenStream.
func NewTestSuiteParser(input antlr.TokenStream) *TestSuiteParser {
	TestSuiteParserInit()
	this := new(TestSuiteParser)
	this.BaseParser = antlr.NewBaseParser(input)
  staticData := &TestSuiteParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "TestSuite.g4"

	return this
}


// TestSuiteParser tokens.
const (
	TestSuiteParserEOF = antlr.TokenEOF
	TestSuiteParserT__0 = 1
	TestSuiteParserT__1 = 2
	TestSuiteParserT__2 = 3
	TestSuiteParserT__3 = 4
	TestSuiteParserT__4 = 5
	TestSuiteParserT__5 = 6
	TestSuiteParserT__6 = 7
	TestSuiteParserT__7 = 8
	TestSuiteParserT__8 = 9
	TestSuiteParserNEGATION = 10
	TestSuiteParserOPERATION = 11
	TestSuiteParserID = 12
	TestSuiteParserSTRING = 13
	TestSuiteParserDID = 14
	TestSuiteParserHEX = 15
	TestSuiteParserHEXDIG = 16
	TestSuiteParserCOMMENT = 17
	TestSuiteParserWS = 18
	TestSuiteParserNL = 19
)

// TestSuiteParser rules.
const (
	TestSuiteParserRULE_relationship_set = 0
	TestSuiteParserRULE_policy_tests = 1
	TestSuiteParserRULE_checks = 2
	TestSuiteParserRULE_check = 3
	TestSuiteParserRULE_implied_relations = 4
	TestSuiteParserRULE_implied_relation = 5
	TestSuiteParserRULE_object_rel = 6
	TestSuiteParserRULE_delegation_assertions = 7
	TestSuiteParserRULE_delegation_assertion = 8
	TestSuiteParserRULE_relationship = 9
	TestSuiteParserRULE_subject = 10
	TestSuiteParserRULE_object = 11
	TestSuiteParserRULE_object_id = 12
	TestSuiteParserRULE_relation = 13
	TestSuiteParserRULE_resource = 14
	TestSuiteParserRULE_actorid = 15
)

// IRelationship_setContext is an interface to support dynamic dispatch.
type IRelationship_setContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRelationship() []IRelationshipContext
	Relationship(i int) IRelationshipContext

	// IsRelationship_setContext differentiates from other interfaces.
	IsRelationship_setContext()
}

type Relationship_setContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRelationship_setContext() *Relationship_setContext {
	var p = new(Relationship_setContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_relationship_set
	return p
}

func InitEmptyRelationship_setContext(p *Relationship_setContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_relationship_set
}

func (*Relationship_setContext) IsRelationship_setContext() {}

func NewRelationship_setContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Relationship_setContext {
	var p = new(Relationship_setContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_relationship_set

	return p
}

func (s *Relationship_setContext) GetParser() antlr.Parser { return s.parser }

func (s *Relationship_setContext) AllRelationship() []IRelationshipContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRelationshipContext); ok {
			len++
		}
	}

	tst := make([]IRelationshipContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRelationshipContext); ok {
			tst[i] = t.(IRelationshipContext)
			i++
		}
	}

	return tst
}

func (s *Relationship_setContext) Relationship(i int) IRelationshipContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationshipContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationshipContext)
}

func (s *Relationship_setContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Relationship_setContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *Relationship_setContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterRelationship_set(s)
	}
}

func (s *Relationship_setContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitRelationship_set(s)
	}
}




func (p *TestSuiteParser) Relationship_set() (localctx IRelationship_setContext) {
	localctx = NewRelationship_setContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, TestSuiteParserRULE_relationship_set)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(33)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	for ok := true; ok; ok = _la == TestSuiteParserID {
		{
			p.SetState(32)
			p.Relationship()
		}


		p.SetState(35)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_la = p.GetTokenStream().LA(1)
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IPolicy_testsContext is an interface to support dynamic dispatch.
type IPolicy_testsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Checks() IChecksContext
	Implied_relations() IImplied_relationsContext
	Delegation_assertions() IDelegation_assertionsContext

	// IsPolicy_testsContext differentiates from other interfaces.
	IsPolicy_testsContext()
}

type Policy_testsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPolicy_testsContext() *Policy_testsContext {
	var p = new(Policy_testsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_policy_tests
	return p
}

func InitEmptyPolicy_testsContext(p *Policy_testsContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_policy_tests
}

func (*Policy_testsContext) IsPolicy_testsContext() {}

func NewPolicy_testsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Policy_testsContext {
	var p = new(Policy_testsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_policy_tests

	return p
}

func (s *Policy_testsContext) GetParser() antlr.Parser { return s.parser }

func (s *Policy_testsContext) Checks() IChecksContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IChecksContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IChecksContext)
}

func (s *Policy_testsContext) Implied_relations() IImplied_relationsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImplied_relationsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImplied_relationsContext)
}

func (s *Policy_testsContext) Delegation_assertions() IDelegation_assertionsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDelegation_assertionsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDelegation_assertionsContext)
}

func (s *Policy_testsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Policy_testsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *Policy_testsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterPolicy_tests(s)
	}
}

func (s *Policy_testsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitPolicy_tests(s)
	}
}




func (p *TestSuiteParser) Policy_tests() (localctx IPolicy_testsContext) {
	localctx = NewPolicy_testsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, TestSuiteParserRULE_policy_tests)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(37)
		p.Checks()
	}
	{
		p.SetState(38)
		p.Implied_relations()
	}
	{
		p.SetState(39)
		p.Delegation_assertions()
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IChecksContext is an interface to support dynamic dispatch.
type IChecksContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllCheck() []ICheckContext
	Check(i int) ICheckContext

	// IsChecksContext differentiates from other interfaces.
	IsChecksContext()
}

type ChecksContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyChecksContext() *ChecksContext {
	var p = new(ChecksContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_checks
	return p
}

func InitEmptyChecksContext(p *ChecksContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_checks
}

func (*ChecksContext) IsChecksContext() {}

func NewChecksContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ChecksContext {
	var p = new(ChecksContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_checks

	return p
}

func (s *ChecksContext) GetParser() antlr.Parser { return s.parser }

func (s *ChecksContext) AllCheck() []ICheckContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICheckContext); ok {
			len++
		}
	}

	tst := make([]ICheckContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICheckContext); ok {
			tst[i] = t.(ICheckContext)
			i++
		}
	}

	return tst
}

func (s *ChecksContext) Check(i int) ICheckContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICheckContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICheckContext)
}

func (s *ChecksContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ChecksContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ChecksContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterChecks(s)
	}
}

func (s *ChecksContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitChecks(s)
	}
}




func (p *TestSuiteParser) Checks() (localctx IChecksContext) {
	localctx = NewChecksContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, TestSuiteParserRULE_checks)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(41)
		p.Match(TestSuiteParserT__0)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(42)
		p.Match(TestSuiteParserT__1)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(46)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	for _la == TestSuiteParserNEGATION || _la == TestSuiteParserID {
		{
			p.SetState(43)
			p.Check()
		}


		p.SetState(48)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(49)
		p.Match(TestSuiteParserT__2)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// ICheckContext is an interface to support dynamic dispatch.
type ICheckContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Relationship() IRelationshipContext
	NEGATION() antlr.TerminalNode

	// IsCheckContext differentiates from other interfaces.
	IsCheckContext()
}

type CheckContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCheckContext() *CheckContext {
	var p = new(CheckContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_check
	return p
}

func InitEmptyCheckContext(p *CheckContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_check
}

func (*CheckContext) IsCheckContext() {}

func NewCheckContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CheckContext {
	var p = new(CheckContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_check

	return p
}

func (s *CheckContext) GetParser() antlr.Parser { return s.parser }

func (s *CheckContext) Relationship() IRelationshipContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationshipContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationshipContext)
}

func (s *CheckContext) NEGATION() antlr.TerminalNode {
	return s.GetToken(TestSuiteParserNEGATION, 0)
}

func (s *CheckContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CheckContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *CheckContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterCheck(s)
	}
}

func (s *CheckContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitCheck(s)
	}
}




func (p *TestSuiteParser) Check() (localctx ICheckContext) {
	localctx = NewCheckContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, TestSuiteParserRULE_check)
	p.SetState(54)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case TestSuiteParserID:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(51)
			p.Relationship()
		}


	case TestSuiteParserNEGATION:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(52)
			p.Match(TestSuiteParserNEGATION)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(53)
			p.Relationship()
		}



	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IImplied_relationsContext is an interface to support dynamic dispatch.
type IImplied_relationsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllImplied_relation() []IImplied_relationContext
	Implied_relation(i int) IImplied_relationContext

	// IsImplied_relationsContext differentiates from other interfaces.
	IsImplied_relationsContext()
}

type Implied_relationsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImplied_relationsContext() *Implied_relationsContext {
	var p = new(Implied_relationsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_implied_relations
	return p
}

func InitEmptyImplied_relationsContext(p *Implied_relationsContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_implied_relations
}

func (*Implied_relationsContext) IsImplied_relationsContext() {}

func NewImplied_relationsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Implied_relationsContext {
	var p = new(Implied_relationsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_implied_relations

	return p
}

func (s *Implied_relationsContext) GetParser() antlr.Parser { return s.parser }

func (s *Implied_relationsContext) AllImplied_relation() []IImplied_relationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IImplied_relationContext); ok {
			len++
		}
	}

	tst := make([]IImplied_relationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IImplied_relationContext); ok {
			tst[i] = t.(IImplied_relationContext)
			i++
		}
	}

	return tst
}

func (s *Implied_relationsContext) Implied_relation(i int) IImplied_relationContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImplied_relationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImplied_relationContext)
}

func (s *Implied_relationsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Implied_relationsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *Implied_relationsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterImplied_relations(s)
	}
}

func (s *Implied_relationsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitImplied_relations(s)
	}
}




func (p *TestSuiteParser) Implied_relations() (localctx IImplied_relationsContext) {
	localctx = NewImplied_relationsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, TestSuiteParserRULE_implied_relations)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(56)
		p.Match(TestSuiteParserT__3)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(57)
		p.Match(TestSuiteParserT__1)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(61)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	for _la == TestSuiteParserID {
		{
			p.SetState(58)
			p.Implied_relation()
		}


		p.SetState(63)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(64)
		p.Match(TestSuiteParserT__2)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IImplied_relationContext is an interface to support dynamic dispatch.
type IImplied_relationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllObject_rel() []IObject_relContext
	Object_rel(i int) IObject_relContext

	// IsImplied_relationContext differentiates from other interfaces.
	IsImplied_relationContext()
}

type Implied_relationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImplied_relationContext() *Implied_relationContext {
	var p = new(Implied_relationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_implied_relation
	return p
}

func InitEmptyImplied_relationContext(p *Implied_relationContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_implied_relation
}

func (*Implied_relationContext) IsImplied_relationContext() {}

func NewImplied_relationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Implied_relationContext {
	var p = new(Implied_relationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_implied_relation

	return p
}

func (s *Implied_relationContext) GetParser() antlr.Parser { return s.parser }

func (s *Implied_relationContext) AllObject_rel() []IObject_relContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IObject_relContext); ok {
			len++
		}
	}

	tst := make([]IObject_relContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IObject_relContext); ok {
			tst[i] = t.(IObject_relContext)
			i++
		}
	}

	return tst
}

func (s *Implied_relationContext) Object_rel(i int) IObject_relContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObject_relContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObject_relContext)
}

func (s *Implied_relationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Implied_relationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *Implied_relationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterImplied_relation(s)
	}
}

func (s *Implied_relationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitImplied_relation(s)
	}
}




func (p *TestSuiteParser) Implied_relation() (localctx IImplied_relationContext) {
	localctx = NewImplied_relationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, TestSuiteParserRULE_implied_relation)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(66)
		p.Object_rel()
	}
	{
		p.SetState(67)
		p.Match(TestSuiteParserT__4)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(68)
		p.Object_rel()
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IObject_relContext is an interface to support dynamic dispatch.
type IObject_relContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Object() IObjectContext
	Relation() IRelationContext

	// IsObject_relContext differentiates from other interfaces.
	IsObject_relContext()
}

type Object_relContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyObject_relContext() *Object_relContext {
	var p = new(Object_relContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_object_rel
	return p
}

func InitEmptyObject_relContext(p *Object_relContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_object_rel
}

func (*Object_relContext) IsObject_relContext() {}

func NewObject_relContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Object_relContext {
	var p = new(Object_relContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_object_rel

	return p
}

func (s *Object_relContext) GetParser() antlr.Parser { return s.parser }

func (s *Object_relContext) Object() IObjectContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectContext)
}

func (s *Object_relContext) Relation() IRelationContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationContext)
}

func (s *Object_relContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Object_relContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *Object_relContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterObject_rel(s)
	}
}

func (s *Object_relContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitObject_rel(s)
	}
}




func (p *TestSuiteParser) Object_rel() (localctx IObject_relContext) {
	localctx = NewObject_relContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, TestSuiteParserRULE_object_rel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(70)
		p.Object()
	}
	{
		p.SetState(71)
		p.Match(TestSuiteParserT__5)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(72)
		p.Relation()
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IDelegation_assertionsContext is an interface to support dynamic dispatch.
type IDelegation_assertionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllDelegation_assertion() []IDelegation_assertionContext
	Delegation_assertion(i int) IDelegation_assertionContext

	// IsDelegation_assertionsContext differentiates from other interfaces.
	IsDelegation_assertionsContext()
}

type Delegation_assertionsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDelegation_assertionsContext() *Delegation_assertionsContext {
	var p = new(Delegation_assertionsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_delegation_assertions
	return p
}

func InitEmptyDelegation_assertionsContext(p *Delegation_assertionsContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_delegation_assertions
}

func (*Delegation_assertionsContext) IsDelegation_assertionsContext() {}

func NewDelegation_assertionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Delegation_assertionsContext {
	var p = new(Delegation_assertionsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_delegation_assertions

	return p
}

func (s *Delegation_assertionsContext) GetParser() antlr.Parser { return s.parser }

func (s *Delegation_assertionsContext) AllDelegation_assertion() []IDelegation_assertionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDelegation_assertionContext); ok {
			len++
		}
	}

	tst := make([]IDelegation_assertionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDelegation_assertionContext); ok {
			tst[i] = t.(IDelegation_assertionContext)
			i++
		}
	}

	return tst
}

func (s *Delegation_assertionsContext) Delegation_assertion(i int) IDelegation_assertionContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDelegation_assertionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDelegation_assertionContext)
}

func (s *Delegation_assertionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Delegation_assertionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *Delegation_assertionsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterDelegation_assertions(s)
	}
}

func (s *Delegation_assertionsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitDelegation_assertions(s)
	}
}




func (p *TestSuiteParser) Delegation_assertions() (localctx IDelegation_assertionsContext) {
	localctx = NewDelegation_assertionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, TestSuiteParserRULE_delegation_assertions)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(74)
		p.Match(TestSuiteParserT__6)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(75)
		p.Match(TestSuiteParserT__1)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	p.SetState(79)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)


	for _la == TestSuiteParserNEGATION || _la == TestSuiteParserDID {
		{
			p.SetState(76)
			p.Delegation_assertion()
		}


		p.SetState(81)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(82)
		p.Match(TestSuiteParserT__2)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IDelegation_assertionContext is an interface to support dynamic dispatch.
type IDelegation_assertionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Actorid() IActoridContext
	OPERATION() antlr.TerminalNode
	Relationship() IRelationshipContext
	NEGATION() antlr.TerminalNode

	// IsDelegation_assertionContext differentiates from other interfaces.
	IsDelegation_assertionContext()
}

type Delegation_assertionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDelegation_assertionContext() *Delegation_assertionContext {
	var p = new(Delegation_assertionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_delegation_assertion
	return p
}

func InitEmptyDelegation_assertionContext(p *Delegation_assertionContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_delegation_assertion
}

func (*Delegation_assertionContext) IsDelegation_assertionContext() {}

func NewDelegation_assertionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Delegation_assertionContext {
	var p = new(Delegation_assertionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_delegation_assertion

	return p
}

func (s *Delegation_assertionContext) GetParser() antlr.Parser { return s.parser }

func (s *Delegation_assertionContext) Actorid() IActoridContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IActoridContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IActoridContext)
}

func (s *Delegation_assertionContext) OPERATION() antlr.TerminalNode {
	return s.GetToken(TestSuiteParserOPERATION, 0)
}

func (s *Delegation_assertionContext) Relationship() IRelationshipContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationshipContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationshipContext)
}

func (s *Delegation_assertionContext) NEGATION() antlr.TerminalNode {
	return s.GetToken(TestSuiteParserNEGATION, 0)
}

func (s *Delegation_assertionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Delegation_assertionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *Delegation_assertionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterDelegation_assertion(s)
	}
}

func (s *Delegation_assertionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitDelegation_assertion(s)
	}
}




func (p *TestSuiteParser) Delegation_assertion() (localctx IDelegation_assertionContext) {
	localctx = NewDelegation_assertionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, TestSuiteParserRULE_delegation_assertion)
	p.SetState(93)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case TestSuiteParserDID:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(84)
			p.Actorid()
		}
		{
			p.SetState(85)
			p.Match(TestSuiteParserOPERATION)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(86)
			p.Relationship()
		}


	case TestSuiteParserNEGATION:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(88)
			p.Match(TestSuiteParserNEGATION)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(89)
			p.Actorid()
		}
		{
			p.SetState(90)
			p.Match(TestSuiteParserOPERATION)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(91)
			p.Relationship()
		}



	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IRelationshipContext is an interface to support dynamic dispatch.
type IRelationshipContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Object() IObjectContext
	Relation() IRelationContext
	Subject() ISubjectContext

	// IsRelationshipContext differentiates from other interfaces.
	IsRelationshipContext()
}

type RelationshipContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRelationshipContext() *RelationshipContext {
	var p = new(RelationshipContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_relationship
	return p
}

func InitEmptyRelationshipContext(p *RelationshipContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_relationship
}

func (*RelationshipContext) IsRelationshipContext() {}

func NewRelationshipContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelationshipContext {
	var p = new(RelationshipContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_relationship

	return p
}

func (s *RelationshipContext) GetParser() antlr.Parser { return s.parser }

func (s *RelationshipContext) Object() IObjectContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectContext)
}

func (s *RelationshipContext) Relation() IRelationContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationContext)
}

func (s *RelationshipContext) Subject() ISubjectContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISubjectContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISubjectContext)
}

func (s *RelationshipContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelationshipContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *RelationshipContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterRelationship(s)
	}
}

func (s *RelationshipContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitRelationship(s)
	}
}




func (p *TestSuiteParser) Relationship() (localctx IRelationshipContext) {
	localctx = NewRelationshipContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, TestSuiteParserRULE_relationship)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(95)
		p.Object()
	}
	{
		p.SetState(96)
		p.Match(TestSuiteParserT__5)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(97)
		p.Relation()
	}
	{
		p.SetState(98)
		p.Match(TestSuiteParserT__7)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(99)
		p.Subject()
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// ISubjectContext is an interface to support dynamic dispatch.
type ISubjectContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsSubjectContext differentiates from other interfaces.
	IsSubjectContext()
}

type SubjectContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySubjectContext() *SubjectContext {
	var p = new(SubjectContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_subject
	return p
}

func InitEmptySubjectContext(p *SubjectContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_subject
}

func (*SubjectContext) IsSubjectContext() {}

func NewSubjectContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubjectContext {
	var p = new(SubjectContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_subject

	return p
}

func (s *SubjectContext) GetParser() antlr.Parser { return s.parser }

func (s *SubjectContext) CopyAll(ctx *SubjectContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *SubjectContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubjectContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}




type Subj_actorContext struct {
	SubjectContext
}

func NewSubj_actorContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *Subj_actorContext {
	var p = new(Subj_actorContext)

	InitEmptySubjectContext(&p.SubjectContext)
	p.parser = parser
	p.CopyAll(ctx.(*SubjectContext))

	return p
}

func (s *Subj_actorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Subj_actorContext) Actorid() IActoridContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IActoridContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IActoridContext)
}


func (s *Subj_actorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterSubj_actor(s)
	}
}

func (s *Subj_actorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitSubj_actor(s)
	}
}


type Subj_usetContext struct {
	SubjectContext
}

func NewSubj_usetContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *Subj_usetContext {
	var p = new(Subj_usetContext)

	InitEmptySubjectContext(&p.SubjectContext)
	p.parser = parser
	p.CopyAll(ctx.(*SubjectContext))

	return p
}

func (s *Subj_usetContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Subj_usetContext) Object() IObjectContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectContext)
}

func (s *Subj_usetContext) Relation() IRelationContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationContext)
}


func (s *Subj_usetContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterSubj_uset(s)
	}
}

func (s *Subj_usetContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitSubj_uset(s)
	}
}


type Subj_objContext struct {
	SubjectContext
}

func NewSubj_objContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *Subj_objContext {
	var p = new(Subj_objContext)

	InitEmptySubjectContext(&p.SubjectContext)
	p.parser = parser
	p.CopyAll(ctx.(*SubjectContext))

	return p
}

func (s *Subj_objContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Subj_objContext) Object() IObjectContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectContext)
}


func (s *Subj_objContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterSubj_obj(s)
	}
}

func (s *Subj_objContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitSubj_obj(s)
	}
}



func (p *TestSuiteParser) Subject() (localctx ISubjectContext) {
	localctx = NewSubjectContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, TestSuiteParserRULE_subject)
	p.SetState(107)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 6, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSubj_usetContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(101)
			p.Object()
		}
		{
			p.SetState(102)
			p.Match(TestSuiteParserT__5)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(103)
			p.Relation()
		}


	case 2:
		localctx = NewSubj_objContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(105)
			p.Object()
		}


	case 3:
		localctx = NewSubj_actorContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(106)
			p.Actorid()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IObjectContext is an interface to support dynamic dispatch.
type IObjectContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Resource() IResourceContext
	Object_id() IObject_idContext

	// IsObjectContext differentiates from other interfaces.
	IsObjectContext()
}

type ObjectContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyObjectContext() *ObjectContext {
	var p = new(ObjectContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_object
	return p
}

func InitEmptyObjectContext(p *ObjectContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_object
}

func (*ObjectContext) IsObjectContext() {}

func NewObjectContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectContext {
	var p = new(ObjectContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_object

	return p
}

func (s *ObjectContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjectContext) Resource() IResourceContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IResourceContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IResourceContext)
}

func (s *ObjectContext) Object_id() IObject_idContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObject_idContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObject_idContext)
}

func (s *ObjectContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjectContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ObjectContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterObject(s)
	}
}

func (s *ObjectContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitObject(s)
	}
}




func (p *TestSuiteParser) Object() (localctx IObjectContext) {
	localctx = NewObjectContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, TestSuiteParserRULE_object)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(109)
		p.Resource()
	}
	{
		p.SetState(110)
		p.Match(TestSuiteParserT__8)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}
	{
		p.SetState(111)
		p.Object_id()
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IObject_idContext is an interface to support dynamic dispatch.
type IObject_idContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsObject_idContext differentiates from other interfaces.
	IsObject_idContext()
}

type Object_idContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyObject_idContext() *Object_idContext {
	var p = new(Object_idContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_object_id
	return p
}

func InitEmptyObject_idContext(p *Object_idContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_object_id
}

func (*Object_idContext) IsObject_idContext() {}

func NewObject_idContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Object_idContext {
	var p = new(Object_idContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_object_id

	return p
}

func (s *Object_idContext) GetParser() antlr.Parser { return s.parser }

func (s *Object_idContext) CopyAll(ctx *Object_idContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *Object_idContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Object_idContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}




type Utf_idContext struct {
	Object_idContext
}

func NewUtf_idContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *Utf_idContext {
	var p = new(Utf_idContext)

	InitEmptyObject_idContext(&p.Object_idContext)
	p.parser = parser
	p.CopyAll(ctx.(*Object_idContext))

	return p
}

func (s *Utf_idContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Utf_idContext) STRING() antlr.TerminalNode {
	return s.GetToken(TestSuiteParserSTRING, 0)
}


func (s *Utf_idContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterUtf_id(s)
	}
}

func (s *Utf_idContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitUtf_id(s)
	}
}


type Ascii_idContext struct {
	Object_idContext
}

func NewAscii_idContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *Ascii_idContext {
	var p = new(Ascii_idContext)

	InitEmptyObject_idContext(&p.Object_idContext)
	p.parser = parser
	p.CopyAll(ctx.(*Object_idContext))

	return p
}

func (s *Ascii_idContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Ascii_idContext) ID() antlr.TerminalNode {
	return s.GetToken(TestSuiteParserID, 0)
}


func (s *Ascii_idContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterAscii_id(s)
	}
}

func (s *Ascii_idContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitAscii_id(s)
	}
}



func (p *TestSuiteParser) Object_id() (localctx IObject_idContext) {
	localctx = NewObject_idContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, TestSuiteParserRULE_object_id)
	p.SetState(115)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case TestSuiteParserID:
		localctx = NewAscii_idContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(113)
			p.Match(TestSuiteParserID)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case TestSuiteParserSTRING:
		localctx = NewUtf_idContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(114)
			p.Match(TestSuiteParserSTRING)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}



	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}


errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IRelationContext is an interface to support dynamic dispatch.
type IRelationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode

	// IsRelationContext differentiates from other interfaces.
	IsRelationContext()
}

type RelationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRelationContext() *RelationContext {
	var p = new(RelationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_relation
	return p
}

func InitEmptyRelationContext(p *RelationContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_relation
}

func (*RelationContext) IsRelationContext() {}

func NewRelationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelationContext {
	var p = new(RelationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_relation

	return p
}

func (s *RelationContext) GetParser() antlr.Parser { return s.parser }

func (s *RelationContext) ID() antlr.TerminalNode {
	return s.GetToken(TestSuiteParserID, 0)
}

func (s *RelationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *RelationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterRelation(s)
	}
}

func (s *RelationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitRelation(s)
	}
}




func (p *TestSuiteParser) Relation() (localctx IRelationContext) {
	localctx = NewRelationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, TestSuiteParserRULE_relation)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(117)
		p.Match(TestSuiteParserID)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IResourceContext is an interface to support dynamic dispatch.
type IResourceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode

	// IsResourceContext differentiates from other interfaces.
	IsResourceContext()
}

type ResourceContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyResourceContext() *ResourceContext {
	var p = new(ResourceContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_resource
	return p
}

func InitEmptyResourceContext(p *ResourceContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_resource
}

func (*ResourceContext) IsResourceContext() {}

func NewResourceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ResourceContext {
	var p = new(ResourceContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_resource

	return p
}

func (s *ResourceContext) GetParser() antlr.Parser { return s.parser }

func (s *ResourceContext) ID() antlr.TerminalNode {
	return s.GetToken(TestSuiteParserID, 0)
}

func (s *ResourceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ResourceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ResourceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterResource(s)
	}
}

func (s *ResourceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitResource(s)
	}
}




func (p *TestSuiteParser) Resource() (localctx IResourceContext) {
	localctx = NewResourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, TestSuiteParserRULE_resource)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(119)
		p.Match(TestSuiteParserID)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// IActoridContext is an interface to support dynamic dispatch.
type IActoridContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DID() antlr.TerminalNode

	// IsActoridContext differentiates from other interfaces.
	IsActoridContext()
}

type ActoridContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyActoridContext() *ActoridContext {
	var p = new(ActoridContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_actorid
	return p
}

func InitEmptyActoridContext(p *ActoridContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TestSuiteParserRULE_actorid
}

func (*ActoridContext) IsActoridContext() {}

func NewActoridContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ActoridContext {
	var p = new(ActoridContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TestSuiteParserRULE_actorid

	return p
}

func (s *ActoridContext) GetParser() antlr.Parser { return s.parser }

func (s *ActoridContext) DID() antlr.TerminalNode {
	return s.GetToken(TestSuiteParserDID, 0)
}

func (s *ActoridContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ActoridContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ActoridContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.EnterActorid(s)
	}
}

func (s *ActoridContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TestSuiteListener); ok {
		listenerT.ExitActorid(s)
	}
}




func (p *TestSuiteParser) Actorid() (localctx IActoridContext) {
	localctx = NewActoridContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, TestSuiteParserRULE_actorid)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(121)
		p.Match(TestSuiteParserDID)
		if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
		}
	}



errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


