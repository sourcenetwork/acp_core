// Code generated from TestSuite.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser

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


type TestSuiteLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames []string
	// TODO: EOF string
}

var TestSuiteLexerLexerStaticData struct {
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

func testsuitelexerLexerInit() {
  staticData := &TestSuiteLexerLexerStaticData
  staticData.ChannelNames = []string{
    "DEFAULT_TOKEN_CHANNEL", "HIDDEN",
  }
  staticData.ModeNames = []string{
    "DEFAULT_MODE",
  }
  staticData.LiteralNames = []string{
    "", "'Checks'", "'{'", "'}'", "'ImpliedRelations'", "'=>'", "'#'", "'DelegationAssertions'", 
    "'@'", "':'", "'!'",
  }
  staticData.SymbolicNames = []string{
    "", "", "", "", "", "", "", "", "", "", "NEGATION", "OPERATION", "ID", 
    "STRING", "DID", "HEX", "HEXDIG", "COMMENT", "WS", "NL",
  }
  staticData.RuleNames = []string{
    "T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "T__6", "T__7", "T__8", 
    "NEGATION", "OPERATION", "ID", "STRING", "DID", "HEX", "HEXDIG", "COMMENT", 
    "WS", "NL",
  }
  staticData.PredictionContextCache = antlr.NewPredictionContextCache()
  staticData.serializedATN = []int32{
	4, 0, 19, 180, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 
	4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 
	10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 
	7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 1, 0, 1, 0, 1, 0, 1, 0, 
	1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 
	1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 
	1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 
	1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 
	1, 6, 1, 6, 1, 7, 1, 7, 1, 8, 1, 8, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 
	10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 3, 10, 112, 
	8, 10, 1, 11, 1, 11, 4, 11, 116, 8, 11, 11, 11, 12, 11, 117, 1, 12, 1, 
	12, 5, 12, 122, 8, 12, 10, 12, 12, 12, 125, 9, 12, 1, 12, 1, 12, 1, 13, 
	1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 4, 13, 135, 8, 13, 11, 13, 12, 13, 136, 
	1, 13, 1, 13, 4, 13, 141, 8, 13, 11, 13, 12, 13, 142, 1, 14, 1, 14, 1, 
	14, 1, 14, 1, 15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 16, 5, 16, 155, 8, 16, 
	10, 16, 12, 16, 158, 9, 16, 1, 16, 3, 16, 161, 8, 16, 1, 16, 1, 16, 1, 
	16, 1, 16, 1, 17, 4, 17, 168, 8, 17, 11, 17, 12, 17, 169, 1, 17, 1, 17, 
	1, 18, 3, 18, 175, 8, 18, 1, 18, 1, 18, 1, 18, 1, 18, 2, 123, 156, 0, 19, 
	1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 11, 
	23, 12, 25, 13, 27, 14, 29, 15, 31, 16, 33, 17, 35, 18, 37, 19, 1, 0, 6, 
	2, 0, 65, 90, 97, 122, 4, 0, 48, 57, 65, 90, 95, 95, 97, 122, 2, 0, 48, 
	57, 97, 122, 5, 0, 45, 46, 48, 57, 65, 90, 95, 95, 97, 122, 3, 0, 48, 57, 
	65, 70, 97, 102, 2, 0, 9, 9, 32, 32, 188, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 
	0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 
	0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 
	0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 
	1, 0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0, 0, 0, 0, 
	35, 1, 0, 0, 0, 0, 37, 1, 0, 0, 0, 1, 39, 1, 0, 0, 0, 3, 46, 1, 0, 0, 0, 
	5, 48, 1, 0, 0, 0, 7, 50, 1, 0, 0, 0, 9, 67, 1, 0, 0, 0, 11, 70, 1, 0, 
	0, 0, 13, 72, 1, 0, 0, 0, 15, 93, 1, 0, 0, 0, 17, 95, 1, 0, 0, 0, 19, 97, 
	1, 0, 0, 0, 21, 111, 1, 0, 0, 0, 23, 113, 1, 0, 0, 0, 25, 119, 1, 0, 0, 
	0, 27, 128, 1, 0, 0, 0, 29, 144, 1, 0, 0, 0, 31, 148, 1, 0, 0, 0, 33, 150, 
	1, 0, 0, 0, 35, 167, 1, 0, 0, 0, 37, 174, 1, 0, 0, 0, 39, 40, 5, 67, 0, 
	0, 40, 41, 5, 104, 0, 0, 41, 42, 5, 101, 0, 0, 42, 43, 5, 99, 0, 0, 43, 
	44, 5, 107, 0, 0, 44, 45, 5, 115, 0, 0, 45, 2, 1, 0, 0, 0, 46, 47, 5, 123, 
	0, 0, 47, 4, 1, 0, 0, 0, 48, 49, 5, 125, 0, 0, 49, 6, 1, 0, 0, 0, 50, 51, 
	5, 73, 0, 0, 51, 52, 5, 109, 0, 0, 52, 53, 5, 112, 0, 0, 53, 54, 5, 108, 
	0, 0, 54, 55, 5, 105, 0, 0, 55, 56, 5, 101, 0, 0, 56, 57, 5, 100, 0, 0, 
	57, 58, 5, 82, 0, 0, 58, 59, 5, 101, 0, 0, 59, 60, 5, 108, 0, 0, 60, 61, 
	5, 97, 0, 0, 61, 62, 5, 116, 0, 0, 62, 63, 5, 105, 0, 0, 63, 64, 5, 111, 
	0, 0, 64, 65, 5, 110, 0, 0, 65, 66, 5, 115, 0, 0, 66, 8, 1, 0, 0, 0, 67, 
	68, 5, 61, 0, 0, 68, 69, 5, 62, 0, 0, 69, 10, 1, 0, 0, 0, 70, 71, 5, 35, 
	0, 0, 71, 12, 1, 0, 0, 0, 72, 73, 5, 68, 0, 0, 73, 74, 5, 101, 0, 0, 74, 
	75, 5, 108, 0, 0, 75, 76, 5, 101, 0, 0, 76, 77, 5, 103, 0, 0, 77, 78, 5, 
	97, 0, 0, 78, 79, 5, 116, 0, 0, 79, 80, 5, 105, 0, 0, 80, 81, 5, 111, 0, 
	0, 81, 82, 5, 110, 0, 0, 82, 83, 5, 65, 0, 0, 83, 84, 5, 115, 0, 0, 84, 
	85, 5, 115, 0, 0, 85, 86, 5, 101, 0, 0, 86, 87, 5, 114, 0, 0, 87, 88, 5, 
	116, 0, 0, 88, 89, 5, 105, 0, 0, 89, 90, 5, 111, 0, 0, 90, 91, 5, 110, 
	0, 0, 91, 92, 5, 115, 0, 0, 92, 14, 1, 0, 0, 0, 93, 94, 5, 64, 0, 0, 94, 
	16, 1, 0, 0, 0, 95, 96, 5, 58, 0, 0, 96, 18, 1, 0, 0, 0, 97, 98, 5, 33, 
	0, 0, 98, 20, 1, 0, 0, 0, 99, 100, 5, 100, 0, 0, 100, 101, 5, 101, 0, 0, 
	101, 102, 5, 108, 0, 0, 102, 103, 5, 101, 0, 0, 103, 104, 5, 116, 0, 0, 
	104, 112, 5, 101, 0, 0, 105, 106, 5, 99, 0, 0, 106, 107, 5, 114, 0, 0, 
	107, 108, 5, 101, 0, 0, 108, 109, 5, 97, 0, 0, 109, 110, 5, 116, 0, 0, 
	110, 112, 5, 101, 0, 0, 111, 99, 1, 0, 0, 0, 111, 105, 1, 0, 0, 0, 112, 
	22, 1, 0, 0, 0, 113, 115, 7, 0, 0, 0, 114, 116, 7, 1, 0, 0, 115, 114, 1, 
	0, 0, 0, 116, 117, 1, 0, 0, 0, 117, 115, 1, 0, 0, 0, 117, 118, 1, 0, 0, 
	0, 118, 24, 1, 0, 0, 0, 119, 123, 5, 34, 0, 0, 120, 122, 9, 0, 0, 0, 121, 
	120, 1, 0, 0, 0, 122, 125, 1, 0, 0, 0, 123, 124, 1, 0, 0, 0, 123, 121, 
	1, 0, 0, 0, 124, 126, 1, 0, 0, 0, 125, 123, 1, 0, 0, 0, 126, 127, 5, 34, 
	0, 0, 127, 26, 1, 0, 0, 0, 128, 129, 5, 100, 0, 0, 129, 130, 5, 105, 0, 
	0, 130, 131, 5, 100, 0, 0, 131, 132, 5, 58, 0, 0, 132, 134, 1, 0, 0, 0, 
	133, 135, 7, 2, 0, 0, 134, 133, 1, 0, 0, 0, 135, 136, 1, 0, 0, 0, 136, 
	134, 1, 0, 0, 0, 136, 137, 1, 0, 0, 0, 137, 138, 1, 0, 0, 0, 138, 140, 
	5, 58, 0, 0, 139, 141, 7, 3, 0, 0, 140, 139, 1, 0, 0, 0, 141, 142, 1, 0, 
	0, 0, 142, 140, 1, 0, 0, 0, 142, 143, 1, 0, 0, 0, 143, 28, 1, 0, 0, 0, 
	144, 145, 5, 37, 0, 0, 145, 146, 3, 31, 15, 0, 146, 147, 3, 31, 15, 0, 
	147, 30, 1, 0, 0, 0, 148, 149, 7, 4, 0, 0, 149, 32, 1, 0, 0, 0, 150, 151, 
	5, 47, 0, 0, 151, 152, 5, 47, 0, 0, 152, 156, 1, 0, 0, 0, 153, 155, 9, 
	0, 0, 0, 154, 153, 1, 0, 0, 0, 155, 158, 1, 0, 0, 0, 156, 157, 1, 0, 0, 
	0, 156, 154, 1, 0, 0, 0, 157, 160, 1, 0, 0, 0, 158, 156, 1, 0, 0, 0, 159, 
	161, 5, 13, 0, 0, 160, 159, 1, 0, 0, 0, 160, 161, 1, 0, 0, 0, 161, 162, 
	1, 0, 0, 0, 162, 163, 5, 10, 0, 0, 163, 164, 1, 0, 0, 0, 164, 165, 6, 16, 
	0, 0, 165, 34, 1, 0, 0, 0, 166, 168, 7, 5, 0, 0, 167, 166, 1, 0, 0, 0, 
	168, 169, 1, 0, 0, 0, 169, 167, 1, 0, 0, 0, 169, 170, 1, 0, 0, 0, 170, 
	171, 1, 0, 0, 0, 171, 172, 6, 17, 0, 0, 172, 36, 1, 0, 0, 0, 173, 175, 
	5, 13, 0, 0, 174, 173, 1, 0, 0, 0, 174, 175, 1, 0, 0, 0, 175, 176, 1, 0, 
	0, 0, 176, 177, 5, 10, 0, 0, 177, 178, 1, 0, 0, 0, 178, 179, 6, 18, 0, 
	0, 179, 38, 1, 0, 0, 0, 10, 0, 111, 117, 123, 136, 142, 156, 160, 169, 
	174, 1, 6, 0, 0,
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

// TestSuiteLexerInit initializes any static state used to implement TestSuiteLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewTestSuiteLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func TestSuiteLexerInit() {
  staticData := &TestSuiteLexerLexerStaticData
  staticData.once.Do(testsuitelexerLexerInit)
}

// NewTestSuiteLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewTestSuiteLexer(input antlr.CharStream) *TestSuiteLexer {
  TestSuiteLexerInit()
	l := new(TestSuiteLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
  staticData := &TestSuiteLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "TestSuite.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// TestSuiteLexer tokens.
const (
	TestSuiteLexerT__0 = 1
	TestSuiteLexerT__1 = 2
	TestSuiteLexerT__2 = 3
	TestSuiteLexerT__3 = 4
	TestSuiteLexerT__4 = 5
	TestSuiteLexerT__5 = 6
	TestSuiteLexerT__6 = 7
	TestSuiteLexerT__7 = 8
	TestSuiteLexerT__8 = 9
	TestSuiteLexerNEGATION = 10
	TestSuiteLexerOPERATION = 11
	TestSuiteLexerID = 12
	TestSuiteLexerSTRING = 13
	TestSuiteLexerDID = 14
	TestSuiteLexerHEX = 15
	TestSuiteLexerHEXDIG = 16
	TestSuiteLexerCOMMENT = 17
	TestSuiteLexerWS = 18
	TestSuiteLexerNL = 19
)

