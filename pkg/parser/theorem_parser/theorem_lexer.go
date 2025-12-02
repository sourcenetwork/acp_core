// Code generated from ./pkg/parser/theorem_parser/Theorem.g4 by ANTLR 4.13.2. DO NOT EDIT.

package theorem_parser
import (
	"fmt"
  	"sync"
	"unicode"
	"github.com/antlr4-go/antlr/v4"
)
// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter


type TheoremLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames []string
	// TODO: EOF string
}

var TheoremLexerLexerStaticData struct {
  once                   sync.Once
  serializedATN          []int32
  ChannelNames           []string
  ModeNames              []string
  LiteralNames           []string
  SymbolicNames          []string
  RuleNames              []string
  PredictionContextCache *antlr.PredictionContextCache
  atn                    *antlr.ATN
  decisionToDFA          []*antlr.DFA
}

func theoremlexerLexerInit() {
  staticData := &TheoremLexerLexerStaticData
  staticData.ChannelNames = []string{
    "DEFAULT_TOKEN_CHANNEL", "HIDDEN",
  }
  staticData.ModeNames = []string{
    "DEFAULT_MODE",
  }
  staticData.LiteralNames = []string{
    "", "'Authorizations'", "'{'", "'}'", "'ImpliedRelations'", "'=>'", 
    "'#'", "'Delegations'", "'>'", "'@'", "':'", "'!'",
  }
  staticData.SymbolicNames = []string{
    "", "", "", "", "", "", "", "", "", "", "", "NEGATION", "ID", "STRING", 
    "DID", "COMMENT", "WS", "NL",
  }
  staticData.RuleNames = []string{
    "T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "T__6", "T__7", "T__8", 
    "T__9", "NEGATION", "ID", "STRING", "DID", "COMMENT", "WS", "NL",
  }
  staticData.PredictionContextCache = antlr.NewPredictionContextCache()
  staticData.serializedATN = []int32{
	4, 0, 17, 157, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 
	4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 
	10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 
	7, 15, 2, 16, 7, 16, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 
	0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 
	3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 
	3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 
	6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 8, 1, 
	8, 1, 9, 1, 9, 1, 10, 1, 10, 1, 11, 1, 11, 4, 11, 99, 8, 11, 11, 11, 12, 
	11, 100, 1, 12, 1, 12, 5, 12, 105, 8, 12, 10, 12, 12, 12, 108, 9, 12, 1, 
	12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 4, 13, 118, 8, 13, 
	11, 13, 12, 13, 119, 1, 13, 1, 13, 4, 13, 124, 8, 13, 11, 13, 12, 13, 125, 
	1, 14, 1, 14, 1, 14, 1, 14, 5, 14, 132, 8, 14, 10, 14, 12, 14, 135, 9, 
	14, 1, 14, 3, 14, 138, 8, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 15, 4, 15, 
	145, 8, 15, 11, 15, 12, 15, 146, 1, 15, 1, 15, 1, 16, 3, 16, 152, 8, 16, 
	1, 16, 1, 16, 1, 16, 1, 16, 2, 106, 133, 0, 17, 1, 1, 3, 2, 5, 3, 7, 4, 
	9, 5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 11, 23, 12, 25, 13, 27, 14, 
	29, 15, 31, 16, 33, 17, 1, 0, 5, 2, 0, 65, 90, 97, 122, 4, 0, 48, 57, 65, 
	90, 95, 95, 97, 122, 2, 0, 48, 57, 97, 122, 5, 0, 45, 46, 48, 57, 65, 90, 
	95, 95, 97, 122, 2, 0, 9, 9, 32, 32, 164, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 
	0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 
	0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 
	0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 
	1, 0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0, 0, 0, 1, 
	35, 1, 0, 0, 0, 3, 50, 1, 0, 0, 0, 5, 52, 1, 0, 0, 0, 7, 54, 1, 0, 0, 0, 
	9, 71, 1, 0, 0, 0, 11, 74, 1, 0, 0, 0, 13, 76, 1, 0, 0, 0, 15, 88, 1, 0, 
	0, 0, 17, 90, 1, 0, 0, 0, 19, 92, 1, 0, 0, 0, 21, 94, 1, 0, 0, 0, 23, 96, 
	1, 0, 0, 0, 25, 102, 1, 0, 0, 0, 27, 111, 1, 0, 0, 0, 29, 127, 1, 0, 0, 
	0, 31, 144, 1, 0, 0, 0, 33, 151, 1, 0, 0, 0, 35, 36, 5, 65, 0, 0, 36, 37, 
	5, 117, 0, 0, 37, 38, 5, 116, 0, 0, 38, 39, 5, 104, 0, 0, 39, 40, 5, 111, 
	0, 0, 40, 41, 5, 114, 0, 0, 41, 42, 5, 105, 0, 0, 42, 43, 5, 122, 0, 0, 
	43, 44, 5, 97, 0, 0, 44, 45, 5, 116, 0, 0, 45, 46, 5, 105, 0, 0, 46, 47, 
	5, 111, 0, 0, 47, 48, 5, 110, 0, 0, 48, 49, 5, 115, 0, 0, 49, 2, 1, 0, 
	0, 0, 50, 51, 5, 123, 0, 0, 51, 4, 1, 0, 0, 0, 52, 53, 5, 125, 0, 0, 53, 
	6, 1, 0, 0, 0, 54, 55, 5, 73, 0, 0, 55, 56, 5, 109, 0, 0, 56, 57, 5, 112, 
	0, 0, 57, 58, 5, 108, 0, 0, 58, 59, 5, 105, 0, 0, 59, 60, 5, 101, 0, 0, 
	60, 61, 5, 100, 0, 0, 61, 62, 5, 82, 0, 0, 62, 63, 5, 101, 0, 0, 63, 64, 
	5, 108, 0, 0, 64, 65, 5, 97, 0, 0, 65, 66, 5, 116, 0, 0, 66, 67, 5, 105, 
	0, 0, 67, 68, 5, 111, 0, 0, 68, 69, 5, 110, 0, 0, 69, 70, 5, 115, 0, 0, 
	70, 8, 1, 0, 0, 0, 71, 72, 5, 61, 0, 0, 72, 73, 5, 62, 0, 0, 73, 10, 1, 
	0, 0, 0, 74, 75, 5, 35, 0, 0, 75, 12, 1, 0, 0, 0, 76, 77, 5, 68, 0, 0, 
	77, 78, 5, 101, 0, 0, 78, 79, 5, 108, 0, 0, 79, 80, 5, 101, 0, 0, 80, 81, 
	5, 103, 0, 0, 81, 82, 5, 97, 0, 0, 82, 83, 5, 116, 0, 0, 83, 84, 5, 105, 
	0, 0, 84, 85, 5, 111, 0, 0, 85, 86, 5, 110, 0, 0, 86, 87, 5, 115, 0, 0, 
	87, 14, 1, 0, 0, 0, 88, 89, 5, 62, 0, 0, 89, 16, 1, 0, 0, 0, 90, 91, 5, 
	64, 0, 0, 91, 18, 1, 0, 0, 0, 92, 93, 5, 58, 0, 0, 93, 20, 1, 0, 0, 0, 
	94, 95, 5, 33, 0, 0, 95, 22, 1, 0, 0, 0, 96, 98, 7, 0, 0, 0, 97, 99, 7, 
	1, 0, 0, 98, 97, 1, 0, 0, 0, 99, 100, 1, 0, 0, 0, 100, 98, 1, 0, 0, 0, 
	100, 101, 1, 0, 0, 0, 101, 24, 1, 0, 0, 0, 102, 106, 5, 34, 0, 0, 103, 
	105, 9, 0, 0, 0, 104, 103, 1, 0, 0, 0, 105, 108, 1, 0, 0, 0, 106, 107, 
	1, 0, 0, 0, 106, 104, 1, 0, 0, 0, 107, 109, 1, 0, 0, 0, 108, 106, 1, 0, 
	0, 0, 109, 110, 5, 34, 0, 0, 110, 26, 1, 0, 0, 0, 111, 112, 5, 100, 0, 
	0, 112, 113, 5, 105, 0, 0, 113, 114, 5, 100, 0, 0, 114, 115, 5, 58, 0, 
	0, 115, 117, 1, 0, 0, 0, 116, 118, 7, 2, 0, 0, 117, 116, 1, 0, 0, 0, 118, 
	119, 1, 0, 0, 0, 119, 117, 1, 0, 0, 0, 119, 120, 1, 0, 0, 0, 120, 121, 
	1, 0, 0, 0, 121, 123, 5, 58, 0, 0, 122, 124, 7, 3, 0, 0, 123, 122, 1, 0, 
	0, 0, 124, 125, 1, 0, 0, 0, 125, 123, 1, 0, 0, 0, 125, 126, 1, 0, 0, 0, 
	126, 28, 1, 0, 0, 0, 127, 128, 5, 47, 0, 0, 128, 129, 5, 47, 0, 0, 129, 
	133, 1, 0, 0, 0, 130, 132, 9, 0, 0, 0, 131, 130, 1, 0, 0, 0, 132, 135, 
	1, 0, 0, 0, 133, 134, 1, 0, 0, 0, 133, 131, 1, 0, 0, 0, 134, 137, 1, 0, 
	0, 0, 135, 133, 1, 0, 0, 0, 136, 138, 5, 13, 0, 0, 137, 136, 1, 0, 0, 0, 
	137, 138, 1, 0, 0, 0, 138, 139, 1, 0, 0, 0, 139, 140, 5, 10, 0, 0, 140, 
	141, 1, 0, 0, 0, 141, 142, 6, 14, 0, 0, 142, 30, 1, 0, 0, 0, 143, 145, 
	7, 4, 0, 0, 144, 143, 1, 0, 0, 0, 145, 146, 1, 0, 0, 0, 146, 144, 1, 0, 
	0, 0, 146, 147, 1, 0, 0, 0, 147, 148, 1, 0, 0, 0, 148, 149, 6, 15, 0, 0, 
	149, 32, 1, 0, 0, 0, 150, 152, 5, 13, 0, 0, 151, 150, 1, 0, 0, 0, 151, 
	152, 1, 0, 0, 0, 152, 153, 1, 0, 0, 0, 153, 154, 5, 10, 0, 0, 154, 155, 
	1, 0, 0, 0, 155, 156, 6, 16, 0, 0, 156, 34, 1, 0, 0, 0, 9, 0, 100, 106, 
	119, 125, 133, 137, 146, 151, 1, 6, 0, 0,
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

// TheoremLexerInit initializes any static state used to implement TheoremLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewTheoremLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func TheoremLexerInit() {
  staticData := &TheoremLexerLexerStaticData
  staticData.once.Do(theoremlexerLexerInit)
}

// NewTheoremLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewTheoremLexer(input antlr.CharStream) *TheoremLexer {
  TheoremLexerInit()
	l := new(TheoremLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
  staticData := &TheoremLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "Theorem.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// TheoremLexer tokens.
const (
	TheoremLexerT__0 = 1
	TheoremLexerT__1 = 2
	TheoremLexerT__2 = 3
	TheoremLexerT__3 = 4
	TheoremLexerT__4 = 5
	TheoremLexerT__5 = 6
	TheoremLexerT__6 = 7
	TheoremLexerT__7 = 8
	TheoremLexerT__8 = 9
	TheoremLexerT__9 = 10
	TheoremLexerNEGATION = 11
	TheoremLexerID = 12
	TheoremLexerSTRING = 13
	TheoremLexerDID = 14
	TheoremLexerCOMMENT = 15
	TheoremLexerWS = 16
	TheoremLexerNL = 17
)

