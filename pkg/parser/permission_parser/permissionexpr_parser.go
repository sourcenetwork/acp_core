// Code generated from ./pkg/parser/permission_parser/PermissionExpr.g4 by ANTLR 4.13.2. DO NOT EDIT.

package permission_parser // PermissionExpr
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


type PermissionExprParser struct {
	*antlr.BaseParser
}

var PermissionExprParserStaticData struct {
  once                   sync.Once
  serializedATN          []int32
  LiteralNames           []string
  SymbolicNames          []string
  RuleNames              []string
  PredictionContextCache *antlr.PredictionContextCache
  atn                    *antlr.ATN
  decisionToDFA          []*antlr.DFA
}

func permissionexprParserInit() {
  staticData := &PermissionExprParserStaticData
  staticData.LiteralNames = []string{
    "", "'->'", "'('", "')'", "'+'", "'-'", "'&'",
  }
  staticData.SymbolicNames = []string{
    "", "", "", "", "", "", "", "IDENTIFIER", "WS",
  }
  staticData.RuleNames = []string{
    "expr", "term", "relation", "resource", "operator",
  }
  staticData.PredictionContextCache = antlr.NewPredictionContextCache()
  staticData.serializedATN = []int32{
	4, 1, 8, 43, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7, 4, 
	1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 5, 0, 18, 8, 0, 10, 0, 12, 0, 
	21, 9, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 32, 
	8, 1, 1, 2, 1, 2, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 3, 4, 41, 8, 4, 1, 4, 0, 
	1, 0, 5, 0, 2, 4, 6, 8, 0, 0, 42, 0, 10, 1, 0, 0, 0, 2, 31, 1, 0, 0, 0, 
	4, 33, 1, 0, 0, 0, 6, 35, 1, 0, 0, 0, 8, 40, 1, 0, 0, 0, 10, 11, 6, 0, 
	-1, 0, 11, 12, 3, 2, 1, 0, 12, 19, 1, 0, 0, 0, 13, 14, 10, 1, 0, 0, 14, 
	15, 3, 8, 4, 0, 15, 16, 3, 2, 1, 0, 16, 18, 1, 0, 0, 0, 17, 13, 1, 0, 0, 
	0, 18, 21, 1, 0, 0, 0, 19, 17, 1, 0, 0, 0, 19, 20, 1, 0, 0, 0, 20, 1, 1, 
	0, 0, 0, 21, 19, 1, 0, 0, 0, 22, 32, 3, 4, 2, 0, 23, 24, 3, 6, 3, 0, 24, 
	25, 5, 1, 0, 0, 25, 26, 3, 4, 2, 0, 26, 32, 1, 0, 0, 0, 27, 28, 5, 2, 0, 
	0, 28, 29, 3, 0, 0, 0, 29, 30, 5, 3, 0, 0, 30, 32, 1, 0, 0, 0, 31, 22, 
	1, 0, 0, 0, 31, 23, 1, 0, 0, 0, 31, 27, 1, 0, 0, 0, 32, 3, 1, 0, 0, 0, 
	33, 34, 5, 7, 0, 0, 34, 5, 1, 0, 0, 0, 35, 36, 5, 7, 0, 0, 36, 7, 1, 0, 
	0, 0, 37, 41, 5, 4, 0, 0, 38, 41, 5, 5, 0, 0, 39, 41, 5, 6, 0, 0, 40, 37, 
	1, 0, 0, 0, 40, 38, 1, 0, 0, 0, 40, 39, 1, 0, 0, 0, 41, 9, 1, 0, 0, 0, 
	3, 19, 31, 40,
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

// PermissionExprParserInit initializes any static state used to implement PermissionExprParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewPermissionExprParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func PermissionExprParserInit() {
  staticData := &PermissionExprParserStaticData
  staticData.once.Do(permissionexprParserInit)
}

// NewPermissionExprParser produces a new parser instance for the optional input antlr.TokenStream.
func NewPermissionExprParser(input antlr.TokenStream) *PermissionExprParser {
	PermissionExprParserInit()
	this := new(PermissionExprParser)
	this.BaseParser = antlr.NewBaseParser(input)
  staticData := &PermissionExprParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "PermissionExpr.g4"

	return this
}


// PermissionExprParser tokens.
const (
	PermissionExprParserEOF = antlr.TokenEOF
	PermissionExprParserT__0 = 1
	PermissionExprParserT__1 = 2
	PermissionExprParserT__2 = 3
	PermissionExprParserT__3 = 4
	PermissionExprParserT__4 = 5
	PermissionExprParserT__5 = 6
	PermissionExprParserIDENTIFIER = 7
	PermissionExprParserWS = 8
)

// PermissionExprParser rules.
const (
	PermissionExprParserRULE_expr = 0
	PermissionExprParserRULE_term = 1
	PermissionExprParserRULE_relation = 2
	PermissionExprParserRULE_resource = 3
	PermissionExprParserRULE_operator = 4
)

// IExprContext is an interface to support dynamic dispatch.
type IExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}

type ExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprContext() *ExprContext {
	var p = new(ExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PermissionExprParserRULE_expr
	return p
}

func InitEmptyExprContext(p *ExprContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PermissionExprParserRULE_expr
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PermissionExprParserRULE_expr

	return p
}

func (s *ExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprContext) CopyAll(ctx *ExprContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}





type AtomContext struct {
	ExprContext
}

func NewAtomContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AtomContext {
	var p = new(AtomContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *AtomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AtomContext) Term() ITermContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITermContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITermContext)
}


func (s *AtomContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PermissionExprVisitor:
		return t.VisitAtom(s)

	default:
		return t.VisitChildren(s)
	}
}


type NestedContext struct {
	ExprContext
}

func NewNestedContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NestedContext {
	var p = new(NestedContext)

	InitEmptyExprContext(&p.ExprContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExprContext))

	return p
}

func (s *NestedContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NestedContext) Expr() IExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *NestedContext) Operator() IOperatorContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOperatorContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOperatorContext)
}

func (s *NestedContext) Term() ITermContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITermContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITermContext)
}


func (s *NestedContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PermissionExprVisitor:
		return t.VisitNested(s)

	default:
		return t.VisitChildren(s)
	}
}



func (p *PermissionExprParser) Expr() (localctx IExprContext) {
	return p.expr(0)
}

func (p *PermissionExprParser) expr(_p int) (localctx IExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 0
	p.EnterRecursionRule(localctx, 0, PermissionExprParserRULE_expr, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	localctx = NewAtomContext(p, localctx)
	p.SetParserRuleContext(localctx)
	_prevctx = localctx

	{
		p.SetState(11)
		p.Term()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(19)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewNestedContext(p, NewExprContext(p, _parentctx, _parentState))
			p.PushNewRecursionContext(localctx, _startState, PermissionExprParserRULE_expr)
			p.SetState(13)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
				goto errorExit
			}
			{
				p.SetState(14)
				p.Operator()
			}
			{
				p.SetState(15)
				p.Term()
			}


		}
		p.SetState(21)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
	    	goto errorExit
	    }
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext())
		if p.HasError() {
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
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}


// ITermContext is an interface to support dynamic dispatch.
type ITermContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsTermContext differentiates from other interfaces.
	IsTermContext()
}

type TermContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTermContext() *TermContext {
	var p = new(TermContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PermissionExprParserRULE_term
	return p
}

func InitEmptyTermContext(p *TermContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PermissionExprParserRULE_term
}

func (*TermContext) IsTermContext() {}

func NewTermContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TermContext {
	var p = new(TermContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PermissionExprParserRULE_term

	return p
}

func (s *TermContext) GetParser() antlr.Parser { return s.parser }

func (s *TermContext) CopyAll(ctx *TermContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *TermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TermContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}




type Ttu_termContext struct {
	TermContext
}

func NewTtu_termContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *Ttu_termContext {
	var p = new(Ttu_termContext)

	InitEmptyTermContext(&p.TermContext)
	p.parser = parser
	p.CopyAll(ctx.(*TermContext))

	return p
}

func (s *Ttu_termContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Ttu_termContext) Resource() IResourceContext {
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

func (s *Ttu_termContext) Relation() IRelationContext {
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


func (s *Ttu_termContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PermissionExprVisitor:
		return t.VisitTtu_term(s)

	default:
		return t.VisitChildren(s)
	}
}


type Expr_termContext struct {
	TermContext
}

func NewExpr_termContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *Expr_termContext {
	var p = new(Expr_termContext)

	InitEmptyTermContext(&p.TermContext)
	p.parser = parser
	p.CopyAll(ctx.(*TermContext))

	return p
}

func (s *Expr_termContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Expr_termContext) Expr() IExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}


func (s *Expr_termContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PermissionExprVisitor:
		return t.VisitExpr_term(s)

	default:
		return t.VisitChildren(s)
	}
}


type Cu_termContext struct {
	TermContext
}

func NewCu_termContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *Cu_termContext {
	var p = new(Cu_termContext)

	InitEmptyTermContext(&p.TermContext)
	p.parser = parser
	p.CopyAll(ctx.(*TermContext))

	return p
}

func (s *Cu_termContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Cu_termContext) Relation() IRelationContext {
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


func (s *Cu_termContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PermissionExprVisitor:
		return t.VisitCu_term(s)

	default:
		return t.VisitChildren(s)
	}
}



func (p *PermissionExprParser) Term() (localctx ITermContext) {
	localctx = NewTermContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, PermissionExprParserRULE_term)
	p.SetState(31)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		localctx = NewCu_termContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(22)
			p.Relation()
		}


	case 2:
		localctx = NewTtu_termContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(23)
			p.Resource()
		}
		{
			p.SetState(24)
			p.Match(PermissionExprParserT__0)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(25)
			p.Relation()
		}


	case 3:
		localctx = NewExpr_termContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(27)
			p.Match(PermissionExprParserT__1)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}
		{
			p.SetState(28)
			p.expr(0)
		}
		{
			p.SetState(29)
			p.Match(PermissionExprParserT__2)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
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


// IRelationContext is an interface to support dynamic dispatch.
type IRelationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode

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
	p.RuleIndex = PermissionExprParserRULE_relation
	return p
}

func InitEmptyRelationContext(p *RelationContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PermissionExprParserRULE_relation
}

func (*RelationContext) IsRelationContext() {}

func NewRelationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelationContext {
	var p = new(RelationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PermissionExprParserRULE_relation

	return p
}

func (s *RelationContext) GetParser() antlr.Parser { return s.parser }

func (s *RelationContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PermissionExprParserIDENTIFIER, 0)
}

func (s *RelationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *RelationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PermissionExprVisitor:
		return t.VisitRelation(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *PermissionExprParser) Relation() (localctx IRelationContext) {
	localctx = NewRelationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, PermissionExprParserRULE_relation)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(33)
		p.Match(PermissionExprParserIDENTIFIER)
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
	IDENTIFIER() antlr.TerminalNode

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
	p.RuleIndex = PermissionExprParserRULE_resource
	return p
}

func InitEmptyResourceContext(p *ResourceContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PermissionExprParserRULE_resource
}

func (*ResourceContext) IsResourceContext() {}

func NewResourceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ResourceContext {
	var p = new(ResourceContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PermissionExprParserRULE_resource

	return p
}

func (s *ResourceContext) GetParser() antlr.Parser { return s.parser }

func (s *ResourceContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PermissionExprParserIDENTIFIER, 0)
}

func (s *ResourceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ResourceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ResourceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PermissionExprVisitor:
		return t.VisitResource(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *PermissionExprParser) Resource() (localctx IResourceContext) {
	localctx = NewResourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, PermissionExprParserRULE_resource)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(35)
		p.Match(PermissionExprParserIDENTIFIER)
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


// IOperatorContext is an interface to support dynamic dispatch.
type IOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsOperatorContext differentiates from other interfaces.
	IsOperatorContext()
}

type OperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOperatorContext() *OperatorContext {
	var p = new(OperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PermissionExprParserRULE_operator
	return p
}

func InitEmptyOperatorContext(p *OperatorContext)  {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PermissionExprParserRULE_operator
}

func (*OperatorContext) IsOperatorContext() {}

func NewOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OperatorContext {
	var p = new(OperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PermissionExprParserRULE_operator

	return p
}

func (s *OperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *OperatorContext) CopyAll(ctx *OperatorContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *OperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}




type IntersectionContext struct {
	OperatorContext
}

func NewIntersectionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IntersectionContext {
	var p = new(IntersectionContext)

	InitEmptyOperatorContext(&p.OperatorContext)
	p.parser = parser
	p.CopyAll(ctx.(*OperatorContext))

	return p
}

func (s *IntersectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntersectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PermissionExprVisitor:
		return t.VisitIntersection(s)

	default:
		return t.VisitChildren(s)
	}
}


type DifferenceContext struct {
	OperatorContext
}

func NewDifferenceContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DifferenceContext {
	var p = new(DifferenceContext)

	InitEmptyOperatorContext(&p.OperatorContext)
	p.parser = parser
	p.CopyAll(ctx.(*OperatorContext))

	return p
}

func (s *DifferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DifferenceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PermissionExprVisitor:
		return t.VisitDifference(s)

	default:
		return t.VisitChildren(s)
	}
}


type UnionContext struct {
	OperatorContext
}

func NewUnionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *UnionContext {
	var p = new(UnionContext)

	InitEmptyOperatorContext(&p.OperatorContext)
	p.parser = parser
	p.CopyAll(ctx.(*OperatorContext))

	return p
}

func (s *UnionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PermissionExprVisitor:
		return t.VisitUnion(s)

	default:
		return t.VisitChildren(s)
	}
}



func (p *PermissionExprParser) Operator() (localctx IOperatorContext) {
	localctx = NewOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, PermissionExprParserRULE_operator)
	p.SetState(40)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case PermissionExprParserT__3:
		localctx = NewUnionContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(37)
			p.Match(PermissionExprParserT__3)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case PermissionExprParserT__4:
		localctx = NewDifferenceContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(38)
			p.Match(PermissionExprParserT__4)
			if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
			}
		}


	case PermissionExprParserT__5:
		localctx = NewIntersectionContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(39)
			p.Match(PermissionExprParserT__5)
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


func (p *PermissionExprParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 0:
			var t *ExprContext = nil
			if localctx != nil { t = localctx.(*ExprContext) }
			return p.Expr_Sempred(t, predIndex)


	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *PermissionExprParser) Expr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
			return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

