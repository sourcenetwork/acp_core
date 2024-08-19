// Code generated from ./internal/parser/Theorem.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // Theorem
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

type TheoremParser struct {
	*antlr.BaseParser
}

var TheoremParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func theoremParserInit() {
	staticData := &TheoremParserStaticData
	staticData.LiteralNames = []string{
		"", "'Authorizations'", "'{'", "'}'", "'ImpliedRelations'", "'=>'",
		"'#'", "'Delegations'", "'>'", "'@'", "':'", "'!'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "NEGATION", "OPERATION",
		"ID", "STRING", "DID", "COMMENT", "WS", "NL",
	}
	staticData.RuleNames = []string{
		"relationship_set", "policy_thorem", "authorization_theorems", "authorization_theorem",
		"implied_relations", "implied_relation", "object_rel", "delegation_theorems",
		"delegation_theorem", "relationship", "operation", "subject", "object",
		"object_id", "relation", "resource", "actorid",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 18, 130, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 1, 0, 5, 0, 36, 8, 0, 10, 0, 12, 0, 39, 9, 0, 1, 1, 1, 1,
		1, 1, 1, 2, 1, 2, 1, 2, 5, 2, 47, 8, 2, 10, 2, 12, 2, 50, 9, 2, 1, 2, 1,
		2, 1, 3, 1, 3, 1, 3, 3, 3, 57, 8, 3, 1, 4, 1, 4, 1, 4, 5, 4, 62, 8, 4,
		10, 4, 12, 4, 65, 9, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6,
		1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 5, 7, 80, 8, 7, 10, 7, 12, 7, 83, 9, 7, 1,
		7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 3, 8, 96,
		8, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 10, 1,
		11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 3, 11, 114, 8, 11, 1, 12, 1, 12,
		1, 12, 1, 12, 1, 13, 1, 13, 3, 13, 122, 8, 13, 1, 14, 1, 14, 1, 15, 1,
		15, 1, 16, 1, 16, 1, 16, 0, 0, 17, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20,
		22, 24, 26, 28, 30, 32, 0, 0, 121, 0, 37, 1, 0, 0, 0, 2, 40, 1, 0, 0, 0,
		4, 43, 1, 0, 0, 0, 6, 56, 1, 0, 0, 0, 8, 58, 1, 0, 0, 0, 10, 68, 1, 0,
		0, 0, 12, 72, 1, 0, 0, 0, 14, 76, 1, 0, 0, 0, 16, 95, 1, 0, 0, 0, 18, 97,
		1, 0, 0, 0, 20, 103, 1, 0, 0, 0, 22, 113, 1, 0, 0, 0, 24, 115, 1, 0, 0,
		0, 26, 121, 1, 0, 0, 0, 28, 123, 1, 0, 0, 0, 30, 125, 1, 0, 0, 0, 32, 127,
		1, 0, 0, 0, 34, 36, 3, 18, 9, 0, 35, 34, 1, 0, 0, 0, 36, 39, 1, 0, 0, 0,
		37, 35, 1, 0, 0, 0, 37, 38, 1, 0, 0, 0, 38, 1, 1, 0, 0, 0, 39, 37, 1, 0,
		0, 0, 40, 41, 3, 4, 2, 0, 41, 42, 3, 14, 7, 0, 42, 3, 1, 0, 0, 0, 43, 44,
		5, 1, 0, 0, 44, 48, 5, 2, 0, 0, 45, 47, 3, 6, 3, 0, 46, 45, 1, 0, 0, 0,
		47, 50, 1, 0, 0, 0, 48, 46, 1, 0, 0, 0, 48, 49, 1, 0, 0, 0, 49, 51, 1,
		0, 0, 0, 50, 48, 1, 0, 0, 0, 51, 52, 5, 3, 0, 0, 52, 5, 1, 0, 0, 0, 53,
		57, 3, 18, 9, 0, 54, 55, 5, 11, 0, 0, 55, 57, 3, 18, 9, 0, 56, 53, 1, 0,
		0, 0, 56, 54, 1, 0, 0, 0, 57, 7, 1, 0, 0, 0, 58, 59, 5, 4, 0, 0, 59, 63,
		5, 2, 0, 0, 60, 62, 3, 10, 5, 0, 61, 60, 1, 0, 0, 0, 62, 65, 1, 0, 0, 0,
		63, 61, 1, 0, 0, 0, 63, 64, 1, 0, 0, 0, 64, 66, 1, 0, 0, 0, 65, 63, 1,
		0, 0, 0, 66, 67, 5, 3, 0, 0, 67, 9, 1, 0, 0, 0, 68, 69, 3, 12, 6, 0, 69,
		70, 5, 5, 0, 0, 70, 71, 3, 12, 6, 0, 71, 11, 1, 0, 0, 0, 72, 73, 3, 24,
		12, 0, 73, 74, 5, 6, 0, 0, 74, 75, 3, 28, 14, 0, 75, 13, 1, 0, 0, 0, 76,
		77, 5, 7, 0, 0, 77, 81, 5, 2, 0, 0, 78, 80, 3, 16, 8, 0, 79, 78, 1, 0,
		0, 0, 80, 83, 1, 0, 0, 0, 81, 79, 1, 0, 0, 0, 81, 82, 1, 0, 0, 0, 82, 84,
		1, 0, 0, 0, 83, 81, 1, 0, 0, 0, 84, 85, 5, 3, 0, 0, 85, 15, 1, 0, 0, 0,
		86, 87, 3, 32, 16, 0, 87, 88, 5, 8, 0, 0, 88, 89, 3, 20, 10, 0, 89, 96,
		1, 0, 0, 0, 90, 91, 5, 11, 0, 0, 91, 92, 3, 32, 16, 0, 92, 93, 5, 8, 0,
		0, 93, 94, 3, 20, 10, 0, 94, 96, 1, 0, 0, 0, 95, 86, 1, 0, 0, 0, 95, 90,
		1, 0, 0, 0, 96, 17, 1, 0, 0, 0, 97, 98, 3, 24, 12, 0, 98, 99, 5, 6, 0,
		0, 99, 100, 3, 28, 14, 0, 100, 101, 5, 9, 0, 0, 101, 102, 3, 22, 11, 0,
		102, 19, 1, 0, 0, 0, 103, 104, 3, 24, 12, 0, 104, 105, 5, 6, 0, 0, 105,
		106, 3, 28, 14, 0, 106, 21, 1, 0, 0, 0, 107, 108, 3, 24, 12, 0, 108, 109,
		5, 6, 0, 0, 109, 110, 3, 28, 14, 0, 110, 114, 1, 0, 0, 0, 111, 114, 3,
		24, 12, 0, 112, 114, 3, 32, 16, 0, 113, 107, 1, 0, 0, 0, 113, 111, 1, 0,
		0, 0, 113, 112, 1, 0, 0, 0, 114, 23, 1, 0, 0, 0, 115, 116, 3, 30, 15, 0,
		116, 117, 5, 10, 0, 0, 117, 118, 3, 26, 13, 0, 118, 25, 1, 0, 0, 0, 119,
		122, 5, 13, 0, 0, 120, 122, 5, 14, 0, 0, 121, 119, 1, 0, 0, 0, 121, 120,
		1, 0, 0, 0, 122, 27, 1, 0, 0, 0, 123, 124, 5, 13, 0, 0, 124, 29, 1, 0,
		0, 0, 125, 126, 5, 13, 0, 0, 126, 31, 1, 0, 0, 0, 127, 128, 5, 15, 0, 0,
		128, 33, 1, 0, 0, 0, 8, 37, 48, 56, 63, 81, 95, 113, 121,
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

// TheoremParserInit initializes any static state used to implement TheoremParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewTheoremParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func TheoremParserInit() {
	staticData := &TheoremParserStaticData
	staticData.once.Do(theoremParserInit)
}

// NewTheoremParser produces a new parser instance for the optional input antlr.TokenStream.
func NewTheoremParser(input antlr.TokenStream) *TheoremParser {
	TheoremParserInit()
	this := new(TheoremParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &TheoremParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "Theorem.g4"

	return this
}

// TheoremParser tokens.
const (
	TheoremParserEOF       = antlr.TokenEOF
	TheoremParserT__0      = 1
	TheoremParserT__1      = 2
	TheoremParserT__2      = 3
	TheoremParserT__3      = 4
	TheoremParserT__4      = 5
	TheoremParserT__5      = 6
	TheoremParserT__6      = 7
	TheoremParserT__7      = 8
	TheoremParserT__8      = 9
	TheoremParserT__9      = 10
	TheoremParserNEGATION  = 11
	TheoremParserOPERATION = 12
	TheoremParserID        = 13
	TheoremParserSTRING    = 14
	TheoremParserDID       = 15
	TheoremParserCOMMENT   = 16
	TheoremParserWS        = 17
	TheoremParserNL        = 18
)

// TheoremParser rules.
const (
	TheoremParserRULE_relationship_set       = 0
	TheoremParserRULE_policy_thorem          = 1
	TheoremParserRULE_authorization_theorems = 2
	TheoremParserRULE_authorization_theorem  = 3
	TheoremParserRULE_implied_relations      = 4
	TheoremParserRULE_implied_relation       = 5
	TheoremParserRULE_object_rel             = 6
	TheoremParserRULE_delegation_theorems    = 7
	TheoremParserRULE_delegation_theorem     = 8
	TheoremParserRULE_relationship           = 9
	TheoremParserRULE_operation              = 10
	TheoremParserRULE_subject                = 11
	TheoremParserRULE_object                 = 12
	TheoremParserRULE_object_id              = 13
	TheoremParserRULE_relation               = 14
	TheoremParserRULE_resource               = 15
	TheoremParserRULE_actorid                = 16
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
	p.RuleIndex = TheoremParserRULE_relationship_set
	return p
}

func InitEmptyRelationship_setContext(p *Relationship_setContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_relationship_set
}

func (*Relationship_setContext) IsRelationship_setContext() {}

func NewRelationship_setContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Relationship_setContext {
	var p = new(Relationship_setContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_relationship_set

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
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationshipContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
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

func (s *Relationship_setContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitRelationship_set(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Relationship_set() (localctx IRelationship_setContext) {
	localctx = NewRelationship_setContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, TheoremParserRULE_relationship_set)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(37)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == TheoremParserID {
		{
			p.SetState(34)
			p.Relationship()
		}

		p.SetState(39)
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

// IPolicy_thoremContext is an interface to support dynamic dispatch.
type IPolicy_thoremContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Authorization_theorems() IAuthorization_theoremsContext
	Delegation_theorems() IDelegation_theoremsContext

	// IsPolicy_thoremContext differentiates from other interfaces.
	IsPolicy_thoremContext()
}

type Policy_thoremContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPolicy_thoremContext() *Policy_thoremContext {
	var p = new(Policy_thoremContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_policy_thorem
	return p
}

func InitEmptyPolicy_thoremContext(p *Policy_thoremContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_policy_thorem
}

func (*Policy_thoremContext) IsPolicy_thoremContext() {}

func NewPolicy_thoremContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Policy_thoremContext {
	var p = new(Policy_thoremContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_policy_thorem

	return p
}

func (s *Policy_thoremContext) GetParser() antlr.Parser { return s.parser }

func (s *Policy_thoremContext) Authorization_theorems() IAuthorization_theoremsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAuthorization_theoremsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAuthorization_theoremsContext)
}

func (s *Policy_thoremContext) Delegation_theorems() IDelegation_theoremsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDelegation_theoremsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDelegation_theoremsContext)
}

func (s *Policy_thoremContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Policy_thoremContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Policy_thoremContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitPolicy_thorem(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Policy_thorem() (localctx IPolicy_thoremContext) {
	localctx = NewPolicy_thoremContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, TheoremParserRULE_policy_thorem)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(40)
		p.Authorization_theorems()
	}
	{
		p.SetState(41)
		p.Delegation_theorems()
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

// IAuthorization_theoremsContext is an interface to support dynamic dispatch.
type IAuthorization_theoremsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAuthorization_theorem() []IAuthorization_theoremContext
	Authorization_theorem(i int) IAuthorization_theoremContext

	// IsAuthorization_theoremsContext differentiates from other interfaces.
	IsAuthorization_theoremsContext()
}

type Authorization_theoremsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAuthorization_theoremsContext() *Authorization_theoremsContext {
	var p = new(Authorization_theoremsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_authorization_theorems
	return p
}

func InitEmptyAuthorization_theoremsContext(p *Authorization_theoremsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_authorization_theorems
}

func (*Authorization_theoremsContext) IsAuthorization_theoremsContext() {}

func NewAuthorization_theoremsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Authorization_theoremsContext {
	var p = new(Authorization_theoremsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_authorization_theorems

	return p
}

func (s *Authorization_theoremsContext) GetParser() antlr.Parser { return s.parser }

func (s *Authorization_theoremsContext) AllAuthorization_theorem() []IAuthorization_theoremContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAuthorization_theoremContext); ok {
			len++
		}
	}

	tst := make([]IAuthorization_theoremContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAuthorization_theoremContext); ok {
			tst[i] = t.(IAuthorization_theoremContext)
			i++
		}
	}

	return tst
}

func (s *Authorization_theoremsContext) Authorization_theorem(i int) IAuthorization_theoremContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAuthorization_theoremContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAuthorization_theoremContext)
}

func (s *Authorization_theoremsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Authorization_theoremsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Authorization_theoremsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitAuthorization_theorems(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Authorization_theorems() (localctx IAuthorization_theoremsContext) {
	localctx = NewAuthorization_theoremsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, TheoremParserRULE_authorization_theorems)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(43)
		p.Match(TheoremParserT__0)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(44)
		p.Match(TheoremParserT__1)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(48)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == TheoremParserNEGATION || _la == TheoremParserID {
		{
			p.SetState(45)
			p.Authorization_theorem()
		}

		p.SetState(50)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(51)
		p.Match(TheoremParserT__2)
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

// IAuthorization_theoremContext is an interface to support dynamic dispatch.
type IAuthorization_theoremContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Relationship() IRelationshipContext
	NEGATION() antlr.TerminalNode

	// IsAuthorization_theoremContext differentiates from other interfaces.
	IsAuthorization_theoremContext()
}

type Authorization_theoremContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAuthorization_theoremContext() *Authorization_theoremContext {
	var p = new(Authorization_theoremContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_authorization_theorem
	return p
}

func InitEmptyAuthorization_theoremContext(p *Authorization_theoremContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_authorization_theorem
}

func (*Authorization_theoremContext) IsAuthorization_theoremContext() {}

func NewAuthorization_theoremContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Authorization_theoremContext {
	var p = new(Authorization_theoremContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_authorization_theorem

	return p
}

func (s *Authorization_theoremContext) GetParser() antlr.Parser { return s.parser }

func (s *Authorization_theoremContext) Relationship() IRelationshipContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationshipContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationshipContext)
}

func (s *Authorization_theoremContext) NEGATION() antlr.TerminalNode {
	return s.GetToken(TheoremParserNEGATION, 0)
}

func (s *Authorization_theoremContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Authorization_theoremContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Authorization_theoremContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitAuthorization_theorem(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Authorization_theorem() (localctx IAuthorization_theoremContext) {
	localctx = NewAuthorization_theoremContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, TheoremParserRULE_authorization_theorem)
	p.SetState(56)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case TheoremParserID:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(53)
			p.Relationship()
		}

	case TheoremParserNEGATION:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(54)
			p.Match(TheoremParserNEGATION)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(55)
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
	p.RuleIndex = TheoremParserRULE_implied_relations
	return p
}

func InitEmptyImplied_relationsContext(p *Implied_relationsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_implied_relations
}

func (*Implied_relationsContext) IsImplied_relationsContext() {}

func NewImplied_relationsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Implied_relationsContext {
	var p = new(Implied_relationsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_implied_relations

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
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImplied_relationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
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

func (s *Implied_relationsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitImplied_relations(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Implied_relations() (localctx IImplied_relationsContext) {
	localctx = NewImplied_relationsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, TheoremParserRULE_implied_relations)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(58)
		p.Match(TheoremParserT__3)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(59)
		p.Match(TheoremParserT__1)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(63)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == TheoremParserID {
		{
			p.SetState(60)
			p.Implied_relation()
		}

		p.SetState(65)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(66)
		p.Match(TheoremParserT__2)
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
	p.RuleIndex = TheoremParserRULE_implied_relation
	return p
}

func InitEmptyImplied_relationContext(p *Implied_relationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_implied_relation
}

func (*Implied_relationContext) IsImplied_relationContext() {}

func NewImplied_relationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Implied_relationContext {
	var p = new(Implied_relationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_implied_relation

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
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObject_relContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
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

func (s *Implied_relationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitImplied_relation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Implied_relation() (localctx IImplied_relationContext) {
	localctx = NewImplied_relationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, TheoremParserRULE_implied_relation)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(68)
		p.Object_rel()
	}
	{
		p.SetState(69)
		p.Match(TheoremParserT__4)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(70)
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
	p.RuleIndex = TheoremParserRULE_object_rel
	return p
}

func InitEmptyObject_relContext(p *Object_relContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_object_rel
}

func (*Object_relContext) IsObject_relContext() {}

func NewObject_relContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Object_relContext {
	var p = new(Object_relContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_object_rel

	return p
}

func (s *Object_relContext) GetParser() antlr.Parser { return s.parser }

func (s *Object_relContext) Object() IObjectContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectContext)
}

func (s *Object_relContext) Relation() IRelationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationContext); ok {
			t = ctx.(antlr.RuleContext)
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

func (s *Object_relContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitObject_rel(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Object_rel() (localctx IObject_relContext) {
	localctx = NewObject_relContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, TheoremParserRULE_object_rel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(72)
		p.Object()
	}
	{
		p.SetState(73)
		p.Match(TheoremParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(74)
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

// IDelegation_theoremsContext is an interface to support dynamic dispatch.
type IDelegation_theoremsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllDelegation_theorem() []IDelegation_theoremContext
	Delegation_theorem(i int) IDelegation_theoremContext

	// IsDelegation_theoremsContext differentiates from other interfaces.
	IsDelegation_theoremsContext()
}

type Delegation_theoremsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDelegation_theoremsContext() *Delegation_theoremsContext {
	var p = new(Delegation_theoremsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_delegation_theorems
	return p
}

func InitEmptyDelegation_theoremsContext(p *Delegation_theoremsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_delegation_theorems
}

func (*Delegation_theoremsContext) IsDelegation_theoremsContext() {}

func NewDelegation_theoremsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Delegation_theoremsContext {
	var p = new(Delegation_theoremsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_delegation_theorems

	return p
}

func (s *Delegation_theoremsContext) GetParser() antlr.Parser { return s.parser }

func (s *Delegation_theoremsContext) AllDelegation_theorem() []IDelegation_theoremContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDelegation_theoremContext); ok {
			len++
		}
	}

	tst := make([]IDelegation_theoremContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDelegation_theoremContext); ok {
			tst[i] = t.(IDelegation_theoremContext)
			i++
		}
	}

	return tst
}

func (s *Delegation_theoremsContext) Delegation_theorem(i int) IDelegation_theoremContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDelegation_theoremContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDelegation_theoremContext)
}

func (s *Delegation_theoremsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Delegation_theoremsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Delegation_theoremsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitDelegation_theorems(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Delegation_theorems() (localctx IDelegation_theoremsContext) {
	localctx = NewDelegation_theoremsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, TheoremParserRULE_delegation_theorems)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(76)
		p.Match(TheoremParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(77)
		p.Match(TheoremParserT__1)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(81)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == TheoremParserNEGATION || _la == TheoremParserDID {
		{
			p.SetState(78)
			p.Delegation_theorem()
		}

		p.SetState(83)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(84)
		p.Match(TheoremParserT__2)
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

// IDelegation_theoremContext is an interface to support dynamic dispatch.
type IDelegation_theoremContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Actorid() IActoridContext
	Operation() IOperationContext
	NEGATION() antlr.TerminalNode

	// IsDelegation_theoremContext differentiates from other interfaces.
	IsDelegation_theoremContext()
}

type Delegation_theoremContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDelegation_theoremContext() *Delegation_theoremContext {
	var p = new(Delegation_theoremContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_delegation_theorem
	return p
}

func InitEmptyDelegation_theoremContext(p *Delegation_theoremContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_delegation_theorem
}

func (*Delegation_theoremContext) IsDelegation_theoremContext() {}

func NewDelegation_theoremContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Delegation_theoremContext {
	var p = new(Delegation_theoremContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_delegation_theorem

	return p
}

func (s *Delegation_theoremContext) GetParser() antlr.Parser { return s.parser }

func (s *Delegation_theoremContext) Actorid() IActoridContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IActoridContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IActoridContext)
}

func (s *Delegation_theoremContext) Operation() IOperationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOperationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOperationContext)
}

func (s *Delegation_theoremContext) NEGATION() antlr.TerminalNode {
	return s.GetToken(TheoremParserNEGATION, 0)
}

func (s *Delegation_theoremContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Delegation_theoremContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Delegation_theoremContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitDelegation_theorem(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Delegation_theorem() (localctx IDelegation_theoremContext) {
	localctx = NewDelegation_theoremContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, TheoremParserRULE_delegation_theorem)
	p.SetState(95)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case TheoremParserDID:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(86)
			p.Actorid()
		}
		{
			p.SetState(87)
			p.Match(TheoremParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(88)
			p.Operation()
		}

	case TheoremParserNEGATION:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(90)
			p.Match(TheoremParserNEGATION)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(91)
			p.Actorid()
		}
		{
			p.SetState(92)
			p.Match(TheoremParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(93)
			p.Operation()
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
	p.RuleIndex = TheoremParserRULE_relationship
	return p
}

func InitEmptyRelationshipContext(p *RelationshipContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_relationship
}

func (*RelationshipContext) IsRelationshipContext() {}

func NewRelationshipContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelationshipContext {
	var p = new(RelationshipContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_relationship

	return p
}

func (s *RelationshipContext) GetParser() antlr.Parser { return s.parser }

func (s *RelationshipContext) Object() IObjectContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectContext)
}

func (s *RelationshipContext) Relation() IRelationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationContext)
}

func (s *RelationshipContext) Subject() ISubjectContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISubjectContext); ok {
			t = ctx.(antlr.RuleContext)
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

func (s *RelationshipContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitRelationship(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Relationship() (localctx IRelationshipContext) {
	localctx = NewRelationshipContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, TheoremParserRULE_relationship)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(97)
		p.Object()
	}
	{
		p.SetState(98)
		p.Match(TheoremParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(99)
		p.Relation()
	}
	{
		p.SetState(100)
		p.Match(TheoremParserT__8)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(101)
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

// IOperationContext is an interface to support dynamic dispatch.
type IOperationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Object() IObjectContext
	Relation() IRelationContext

	// IsOperationContext differentiates from other interfaces.
	IsOperationContext()
}

type OperationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOperationContext() *OperationContext {
	var p = new(OperationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_operation
	return p
}

func InitEmptyOperationContext(p *OperationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_operation
}

func (*OperationContext) IsOperationContext() {}

func NewOperationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OperationContext {
	var p = new(OperationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_operation

	return p
}

func (s *OperationContext) GetParser() antlr.Parser { return s.parser }

func (s *OperationContext) Object() IObjectContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectContext)
}

func (s *OperationContext) Relation() IRelationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationContext)
}

func (s *OperationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OperationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OperationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitOperation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Operation() (localctx IOperationContext) {
	localctx = NewOperationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, TheoremParserRULE_operation)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(103)
		p.Object()
	}
	{
		p.SetState(104)
		p.Match(TheoremParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(105)
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
	p.RuleIndex = TheoremParserRULE_subject
	return p
}

func InitEmptySubjectContext(p *SubjectContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_subject
}

func (*SubjectContext) IsSubjectContext() {}

func NewSubjectContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubjectContext {
	var p = new(SubjectContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_subject

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
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IActoridContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IActoridContext)
}

func (s *Subj_actorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitSubj_actor(s)

	default:
		return t.VisitChildren(s)
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
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectContext)
}

func (s *Subj_usetContext) Relation() IRelationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationContext)
}

func (s *Subj_usetContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitSubj_uset(s)

	default:
		return t.VisitChildren(s)
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
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectContext)
}

func (s *Subj_objContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitSubj_obj(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Subject() (localctx ISubjectContext) {
	localctx = NewSubjectContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, TheoremParserRULE_subject)
	p.SetState(113)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 6, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSubj_usetContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(107)
			p.Object()
		}
		{
			p.SetState(108)
			p.Match(TheoremParserT__5)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(109)
			p.Relation()
		}

	case 2:
		localctx = NewSubj_objContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(111)
			p.Object()
		}

	case 3:
		localctx = NewSubj_actorContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(112)
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
	p.RuleIndex = TheoremParserRULE_object
	return p
}

func InitEmptyObjectContext(p *ObjectContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_object
}

func (*ObjectContext) IsObjectContext() {}

func NewObjectContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectContext {
	var p = new(ObjectContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_object

	return p
}

func (s *ObjectContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjectContext) Resource() IResourceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IResourceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IResourceContext)
}

func (s *ObjectContext) Object_id() IObject_idContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObject_idContext); ok {
			t = ctx.(antlr.RuleContext)
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

func (s *ObjectContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitObject(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Object() (localctx IObjectContext) {
	localctx = NewObjectContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, TheoremParserRULE_object)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(115)
		p.Resource()
	}
	{
		p.SetState(116)
		p.Match(TheoremParserT__9)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(117)
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
	p.RuleIndex = TheoremParserRULE_object_id
	return p
}

func InitEmptyObject_idContext(p *Object_idContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_object_id
}

func (*Object_idContext) IsObject_idContext() {}

func NewObject_idContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Object_idContext {
	var p = new(Object_idContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_object_id

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
	return s.GetToken(TheoremParserSTRING, 0)
}

func (s *Utf_idContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitUtf_id(s)

	default:
		return t.VisitChildren(s)
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
	return s.GetToken(TheoremParserID, 0)
}

func (s *Ascii_idContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitAscii_id(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Object_id() (localctx IObject_idContext) {
	localctx = NewObject_idContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, TheoremParserRULE_object_id)
	p.SetState(121)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case TheoremParserID:
		localctx = NewAscii_idContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(119)
			p.Match(TheoremParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case TheoremParserSTRING:
		localctx = NewUtf_idContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(120)
			p.Match(TheoremParserSTRING)
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
	p.RuleIndex = TheoremParserRULE_relation
	return p
}

func InitEmptyRelationContext(p *RelationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_relation
}

func (*RelationContext) IsRelationContext() {}

func NewRelationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelationContext {
	var p = new(RelationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_relation

	return p
}

func (s *RelationContext) GetParser() antlr.Parser { return s.parser }

func (s *RelationContext) ID() antlr.TerminalNode {
	return s.GetToken(TheoremParserID, 0)
}

func (s *RelationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RelationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitRelation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Relation() (localctx IRelationContext) {
	localctx = NewRelationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, TheoremParserRULE_relation)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(123)
		p.Match(TheoremParserID)
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
	p.RuleIndex = TheoremParserRULE_resource
	return p
}

func InitEmptyResourceContext(p *ResourceContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_resource
}

func (*ResourceContext) IsResourceContext() {}

func NewResourceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ResourceContext {
	var p = new(ResourceContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_resource

	return p
}

func (s *ResourceContext) GetParser() antlr.Parser { return s.parser }

func (s *ResourceContext) ID() antlr.TerminalNode {
	return s.GetToken(TheoremParserID, 0)
}

func (s *ResourceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ResourceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ResourceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitResource(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Resource() (localctx IResourceContext) {
	localctx = NewResourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, TheoremParserRULE_resource)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(125)
		p.Match(TheoremParserID)
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
	p.RuleIndex = TheoremParserRULE_actorid
	return p
}

func InitEmptyActoridContext(p *ActoridContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TheoremParserRULE_actorid
}

func (*ActoridContext) IsActoridContext() {}

func NewActoridContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ActoridContext {
	var p = new(ActoridContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TheoremParserRULE_actorid

	return p
}

func (s *ActoridContext) GetParser() antlr.Parser { return s.parser }

func (s *ActoridContext) DID() antlr.TerminalNode {
	return s.GetToken(TheoremParserDID, 0)
}

func (s *ActoridContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ActoridContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ActoridContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TheoremVisitor:
		return t.VisitActorid(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TheoremParser) Actorid() (localctx IActoridContext) {
	localctx = NewActoridContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, TheoremParserRULE_actorid)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(127)
		p.Match(TheoremParserDID)
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
