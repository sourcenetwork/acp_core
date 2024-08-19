// Code generated from ./internal/parser/Theorem.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type TheoremLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
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
		"", "", "", "", "", "", "", "", "", "", "", "NEGATION", "OPERATION",
		"ID", "STRING", "DID", "COMMENT", "WS", "NL",
	}
	staticData.RuleNames = []string{
		"T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "T__6", "T__7", "T__8",
		"T__9", "NEGATION", "OPERATION", "ID", "STRING", "DID", "COMMENT", "WS",
		"NL",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 18, 173, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
		1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 2,
		1, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3,
		1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 6,
		1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7,
		1, 7, 1, 8, 1, 8, 1, 9, 1, 9, 1, 10, 1, 10, 1, 11, 1, 11, 1, 11, 1, 11,
		1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 3, 11, 111, 8,
		11, 1, 12, 1, 12, 4, 12, 115, 8, 12, 11, 12, 12, 12, 116, 1, 13, 1, 13,
		5, 13, 121, 8, 13, 10, 13, 12, 13, 124, 9, 13, 1, 13, 1, 13, 1, 14, 1,
		14, 1, 14, 1, 14, 1, 14, 1, 14, 4, 14, 134, 8, 14, 11, 14, 12, 14, 135,
		1, 14, 1, 14, 4, 14, 140, 8, 14, 11, 14, 12, 14, 141, 1, 15, 1, 15, 1,
		15, 1, 15, 5, 15, 148, 8, 15, 10, 15, 12, 15, 151, 9, 15, 1, 15, 3, 15,
		154, 8, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 16, 4, 16, 161, 8, 16, 11, 16,
		12, 16, 162, 1, 16, 1, 16, 1, 17, 3, 17, 168, 8, 17, 1, 17, 1, 17, 1, 17,
		1, 17, 2, 122, 149, 0, 18, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7,
		15, 8, 17, 9, 19, 10, 21, 11, 23, 12, 25, 13, 27, 14, 29, 15, 31, 16, 33,
		17, 35, 18, 1, 0, 5, 2, 0, 65, 90, 97, 122, 4, 0, 48, 57, 65, 90, 95, 95,
		97, 122, 2, 0, 48, 57, 97, 122, 5, 0, 45, 46, 48, 57, 65, 90, 95, 95, 97,
		122, 2, 0, 9, 9, 32, 32, 181, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5,
		1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13,
		1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0,
		21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 1, 0, 0, 0,
		0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0, 0, 0, 0, 35, 1, 0, 0,
		0, 1, 37, 1, 0, 0, 0, 3, 52, 1, 0, 0, 0, 5, 54, 1, 0, 0, 0, 7, 56, 1, 0,
		0, 0, 9, 73, 1, 0, 0, 0, 11, 76, 1, 0, 0, 0, 13, 78, 1, 0, 0, 0, 15, 90,
		1, 0, 0, 0, 17, 92, 1, 0, 0, 0, 19, 94, 1, 0, 0, 0, 21, 96, 1, 0, 0, 0,
		23, 110, 1, 0, 0, 0, 25, 112, 1, 0, 0, 0, 27, 118, 1, 0, 0, 0, 29, 127,
		1, 0, 0, 0, 31, 143, 1, 0, 0, 0, 33, 160, 1, 0, 0, 0, 35, 167, 1, 0, 0,
		0, 37, 38, 5, 65, 0, 0, 38, 39, 5, 117, 0, 0, 39, 40, 5, 116, 0, 0, 40,
		41, 5, 104, 0, 0, 41, 42, 5, 111, 0, 0, 42, 43, 5, 114, 0, 0, 43, 44, 5,
		105, 0, 0, 44, 45, 5, 122, 0, 0, 45, 46, 5, 97, 0, 0, 46, 47, 5, 116, 0,
		0, 47, 48, 5, 105, 0, 0, 48, 49, 5, 111, 0, 0, 49, 50, 5, 110, 0, 0, 50,
		51, 5, 115, 0, 0, 51, 2, 1, 0, 0, 0, 52, 53, 5, 123, 0, 0, 53, 4, 1, 0,
		0, 0, 54, 55, 5, 125, 0, 0, 55, 6, 1, 0, 0, 0, 56, 57, 5, 73, 0, 0, 57,
		58, 5, 109, 0, 0, 58, 59, 5, 112, 0, 0, 59, 60, 5, 108, 0, 0, 60, 61, 5,
		105, 0, 0, 61, 62, 5, 101, 0, 0, 62, 63, 5, 100, 0, 0, 63, 64, 5, 82, 0,
		0, 64, 65, 5, 101, 0, 0, 65, 66, 5, 108, 0, 0, 66, 67, 5, 97, 0, 0, 67,
		68, 5, 116, 0, 0, 68, 69, 5, 105, 0, 0, 69, 70, 5, 111, 0, 0, 70, 71, 5,
		110, 0, 0, 71, 72, 5, 115, 0, 0, 72, 8, 1, 0, 0, 0, 73, 74, 5, 61, 0, 0,
		74, 75, 5, 62, 0, 0, 75, 10, 1, 0, 0, 0, 76, 77, 5, 35, 0, 0, 77, 12, 1,
		0, 0, 0, 78, 79, 5, 68, 0, 0, 79, 80, 5, 101, 0, 0, 80, 81, 5, 108, 0,
		0, 81, 82, 5, 101, 0, 0, 82, 83, 5, 103, 0, 0, 83, 84, 5, 97, 0, 0, 84,
		85, 5, 116, 0, 0, 85, 86, 5, 105, 0, 0, 86, 87, 5, 111, 0, 0, 87, 88, 5,
		110, 0, 0, 88, 89, 5, 115, 0, 0, 89, 14, 1, 0, 0, 0, 90, 91, 5, 62, 0,
		0, 91, 16, 1, 0, 0, 0, 92, 93, 5, 64, 0, 0, 93, 18, 1, 0, 0, 0, 94, 95,
		5, 58, 0, 0, 95, 20, 1, 0, 0, 0, 96, 97, 5, 33, 0, 0, 97, 22, 1, 0, 0,
		0, 98, 99, 5, 100, 0, 0, 99, 100, 5, 101, 0, 0, 100, 101, 5, 108, 0, 0,
		101, 102, 5, 101, 0, 0, 102, 103, 5, 116, 0, 0, 103, 111, 5, 101, 0, 0,
		104, 105, 5, 99, 0, 0, 105, 106, 5, 114, 0, 0, 106, 107, 5, 101, 0, 0,
		107, 108, 5, 97, 0, 0, 108, 109, 5, 116, 0, 0, 109, 111, 5, 101, 0, 0,
		110, 98, 1, 0, 0, 0, 110, 104, 1, 0, 0, 0, 111, 24, 1, 0, 0, 0, 112, 114,
		7, 0, 0, 0, 113, 115, 7, 1, 0, 0, 114, 113, 1, 0, 0, 0, 115, 116, 1, 0,
		0, 0, 116, 114, 1, 0, 0, 0, 116, 117, 1, 0, 0, 0, 117, 26, 1, 0, 0, 0,
		118, 122, 5, 34, 0, 0, 119, 121, 9, 0, 0, 0, 120, 119, 1, 0, 0, 0, 121,
		124, 1, 0, 0, 0, 122, 123, 1, 0, 0, 0, 122, 120, 1, 0, 0, 0, 123, 125,
		1, 0, 0, 0, 124, 122, 1, 0, 0, 0, 125, 126, 5, 34, 0, 0, 126, 28, 1, 0,
		0, 0, 127, 128, 5, 100, 0, 0, 128, 129, 5, 105, 0, 0, 129, 130, 5, 100,
		0, 0, 130, 131, 5, 58, 0, 0, 131, 133, 1, 0, 0, 0, 132, 134, 7, 2, 0, 0,
		133, 132, 1, 0, 0, 0, 134, 135, 1, 0, 0, 0, 135, 133, 1, 0, 0, 0, 135,
		136, 1, 0, 0, 0, 136, 137, 1, 0, 0, 0, 137, 139, 5, 58, 0, 0, 138, 140,
		7, 3, 0, 0, 139, 138, 1, 0, 0, 0, 140, 141, 1, 0, 0, 0, 141, 139, 1, 0,
		0, 0, 141, 142, 1, 0, 0, 0, 142, 30, 1, 0, 0, 0, 143, 144, 5, 47, 0, 0,
		144, 145, 5, 47, 0, 0, 145, 149, 1, 0, 0, 0, 146, 148, 9, 0, 0, 0, 147,
		146, 1, 0, 0, 0, 148, 151, 1, 0, 0, 0, 149, 150, 1, 0, 0, 0, 149, 147,
		1, 0, 0, 0, 150, 153, 1, 0, 0, 0, 151, 149, 1, 0, 0, 0, 152, 154, 5, 13,
		0, 0, 153, 152, 1, 0, 0, 0, 153, 154, 1, 0, 0, 0, 154, 155, 1, 0, 0, 0,
		155, 156, 5, 10, 0, 0, 156, 157, 1, 0, 0, 0, 157, 158, 6, 15, 0, 0, 158,
		32, 1, 0, 0, 0, 159, 161, 7, 4, 0, 0, 160, 159, 1, 0, 0, 0, 161, 162, 1,
		0, 0, 0, 162, 160, 1, 0, 0, 0, 162, 163, 1, 0, 0, 0, 163, 164, 1, 0, 0,
		0, 164, 165, 6, 16, 0, 0, 165, 34, 1, 0, 0, 0, 166, 168, 5, 13, 0, 0, 167,
		166, 1, 0, 0, 0, 167, 168, 1, 0, 0, 0, 168, 169, 1, 0, 0, 0, 169, 170,
		5, 10, 0, 0, 170, 171, 1, 0, 0, 0, 171, 172, 6, 17, 0, 0, 172, 36, 1, 0,
		0, 0, 10, 0, 110, 116, 122, 135, 141, 149, 153, 162, 167, 1, 6, 0, 0,
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
	TheoremLexerT__0      = 1
	TheoremLexerT__1      = 2
	TheoremLexerT__2      = 3
	TheoremLexerT__3      = 4
	TheoremLexerT__4      = 5
	TheoremLexerT__5      = 6
	TheoremLexerT__6      = 7
	TheoremLexerT__7      = 8
	TheoremLexerT__8      = 9
	TheoremLexerT__9      = 10
	TheoremLexerNEGATION  = 11
	TheoremLexerOPERATION = 12
	TheoremLexerID        = 13
	TheoremLexerSTRING    = 14
	TheoremLexerDID       = 15
	TheoremLexerCOMMENT   = 16
	TheoremLexerWS        = 17
	TheoremLexerNL        = 18
)
