// Code generated from Theorem.g4 by ANTLR 4.13.1. DO NOT EDIT.

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
		"ID", "STRING", "DID", "HEX", "HEXDIG", "COMMENT", "WS", "NL",
	}
	staticData.RuleNames = []string{
		"T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "T__6", "T__7", "T__8",
		"T__9", "NEGATION", "OPERATION", "ID", "STRING", "DID", "HEX", "HEXDIG",
		"COMMENT", "WS", "NL",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 20, 183, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 1, 0, 1,
		0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1,
		3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1,
		4, 1, 4, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1,
		6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 8, 1, 8, 1, 9, 1, 9, 1, 10, 1, 10,
		1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1,
		11, 1, 11, 3, 11, 115, 8, 11, 1, 12, 1, 12, 4, 12, 119, 8, 12, 11, 12,
		12, 12, 120, 1, 13, 1, 13, 5, 13, 125, 8, 13, 10, 13, 12, 13, 128, 9, 13,
		1, 13, 1, 13, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 4, 14, 138, 8,
		14, 11, 14, 12, 14, 139, 1, 14, 1, 14, 4, 14, 144, 8, 14, 11, 14, 12, 14,
		145, 1, 15, 1, 15, 1, 15, 1, 15, 1, 16, 1, 16, 1, 17, 1, 17, 1, 17, 1,
		17, 5, 17, 158, 8, 17, 10, 17, 12, 17, 161, 9, 17, 1, 17, 3, 17, 164, 8,
		17, 1, 17, 1, 17, 1, 17, 1, 17, 1, 18, 4, 18, 171, 8, 18, 11, 18, 12, 18,
		172, 1, 18, 1, 18, 1, 19, 3, 19, 178, 8, 19, 1, 19, 1, 19, 1, 19, 1, 19,
		2, 126, 159, 0, 20, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8,
		17, 9, 19, 10, 21, 11, 23, 12, 25, 13, 27, 14, 29, 15, 31, 16, 33, 17,
		35, 18, 37, 19, 39, 20, 1, 0, 6, 2, 0, 65, 90, 97, 122, 4, 0, 48, 57, 65,
		90, 95, 95, 97, 122, 2, 0, 48, 57, 97, 122, 5, 0, 45, 46, 48, 57, 65, 90,
		95, 95, 97, 122, 3, 0, 48, 57, 65, 70, 97, 102, 2, 0, 9, 9, 32, 32, 191,
		0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0,
		0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0,
		0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0,
		0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 1, 0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1,
		0, 0, 0, 0, 33, 1, 0, 0, 0, 0, 35, 1, 0, 0, 0, 0, 37, 1, 0, 0, 0, 0, 39,
		1, 0, 0, 0, 1, 41, 1, 0, 0, 0, 3, 56, 1, 0, 0, 0, 5, 58, 1, 0, 0, 0, 7,
		60, 1, 0, 0, 0, 9, 77, 1, 0, 0, 0, 11, 80, 1, 0, 0, 0, 13, 82, 1, 0, 0,
		0, 15, 94, 1, 0, 0, 0, 17, 96, 1, 0, 0, 0, 19, 98, 1, 0, 0, 0, 21, 100,
		1, 0, 0, 0, 23, 114, 1, 0, 0, 0, 25, 116, 1, 0, 0, 0, 27, 122, 1, 0, 0,
		0, 29, 131, 1, 0, 0, 0, 31, 147, 1, 0, 0, 0, 33, 151, 1, 0, 0, 0, 35, 153,
		1, 0, 0, 0, 37, 170, 1, 0, 0, 0, 39, 177, 1, 0, 0, 0, 41, 42, 5, 65, 0,
		0, 42, 43, 5, 117, 0, 0, 43, 44, 5, 116, 0, 0, 44, 45, 5, 104, 0, 0, 45,
		46, 5, 111, 0, 0, 46, 47, 5, 114, 0, 0, 47, 48, 5, 105, 0, 0, 48, 49, 5,
		122, 0, 0, 49, 50, 5, 97, 0, 0, 50, 51, 5, 116, 0, 0, 51, 52, 5, 105, 0,
		0, 52, 53, 5, 111, 0, 0, 53, 54, 5, 110, 0, 0, 54, 55, 5, 115, 0, 0, 55,
		2, 1, 0, 0, 0, 56, 57, 5, 123, 0, 0, 57, 4, 1, 0, 0, 0, 58, 59, 5, 125,
		0, 0, 59, 6, 1, 0, 0, 0, 60, 61, 5, 73, 0, 0, 61, 62, 5, 109, 0, 0, 62,
		63, 5, 112, 0, 0, 63, 64, 5, 108, 0, 0, 64, 65, 5, 105, 0, 0, 65, 66, 5,
		101, 0, 0, 66, 67, 5, 100, 0, 0, 67, 68, 5, 82, 0, 0, 68, 69, 5, 101, 0,
		0, 69, 70, 5, 108, 0, 0, 70, 71, 5, 97, 0, 0, 71, 72, 5, 116, 0, 0, 72,
		73, 5, 105, 0, 0, 73, 74, 5, 111, 0, 0, 74, 75, 5, 110, 0, 0, 75, 76, 5,
		115, 0, 0, 76, 8, 1, 0, 0, 0, 77, 78, 5, 61, 0, 0, 78, 79, 5, 62, 0, 0,
		79, 10, 1, 0, 0, 0, 80, 81, 5, 35, 0, 0, 81, 12, 1, 0, 0, 0, 82, 83, 5,
		68, 0, 0, 83, 84, 5, 101, 0, 0, 84, 85, 5, 108, 0, 0, 85, 86, 5, 101, 0,
		0, 86, 87, 5, 103, 0, 0, 87, 88, 5, 97, 0, 0, 88, 89, 5, 116, 0, 0, 89,
		90, 5, 105, 0, 0, 90, 91, 5, 111, 0, 0, 91, 92, 5, 110, 0, 0, 92, 93, 5,
		115, 0, 0, 93, 14, 1, 0, 0, 0, 94, 95, 5, 62, 0, 0, 95, 16, 1, 0, 0, 0,
		96, 97, 5, 64, 0, 0, 97, 18, 1, 0, 0, 0, 98, 99, 5, 58, 0, 0, 99, 20, 1,
		0, 0, 0, 100, 101, 5, 33, 0, 0, 101, 22, 1, 0, 0, 0, 102, 103, 5, 100,
		0, 0, 103, 104, 5, 101, 0, 0, 104, 105, 5, 108, 0, 0, 105, 106, 5, 101,
		0, 0, 106, 107, 5, 116, 0, 0, 107, 115, 5, 101, 0, 0, 108, 109, 5, 99,
		0, 0, 109, 110, 5, 114, 0, 0, 110, 111, 5, 101, 0, 0, 111, 112, 5, 97,
		0, 0, 112, 113, 5, 116, 0, 0, 113, 115, 5, 101, 0, 0, 114, 102, 1, 0, 0,
		0, 114, 108, 1, 0, 0, 0, 115, 24, 1, 0, 0, 0, 116, 118, 7, 0, 0, 0, 117,
		119, 7, 1, 0, 0, 118, 117, 1, 0, 0, 0, 119, 120, 1, 0, 0, 0, 120, 118,
		1, 0, 0, 0, 120, 121, 1, 0, 0, 0, 121, 26, 1, 0, 0, 0, 122, 126, 5, 34,
		0, 0, 123, 125, 9, 0, 0, 0, 124, 123, 1, 0, 0, 0, 125, 128, 1, 0, 0, 0,
		126, 127, 1, 0, 0, 0, 126, 124, 1, 0, 0, 0, 127, 129, 1, 0, 0, 0, 128,
		126, 1, 0, 0, 0, 129, 130, 5, 34, 0, 0, 130, 28, 1, 0, 0, 0, 131, 132,
		5, 100, 0, 0, 132, 133, 5, 105, 0, 0, 133, 134, 5, 100, 0, 0, 134, 135,
		5, 58, 0, 0, 135, 137, 1, 0, 0, 0, 136, 138, 7, 2, 0, 0, 137, 136, 1, 0,
		0, 0, 138, 139, 1, 0, 0, 0, 139, 137, 1, 0, 0, 0, 139, 140, 1, 0, 0, 0,
		140, 141, 1, 0, 0, 0, 141, 143, 5, 58, 0, 0, 142, 144, 7, 3, 0, 0, 143,
		142, 1, 0, 0, 0, 144, 145, 1, 0, 0, 0, 145, 143, 1, 0, 0, 0, 145, 146,
		1, 0, 0, 0, 146, 30, 1, 0, 0, 0, 147, 148, 5, 37, 0, 0, 148, 149, 3, 33,
		16, 0, 149, 150, 3, 33, 16, 0, 150, 32, 1, 0, 0, 0, 151, 152, 7, 4, 0,
		0, 152, 34, 1, 0, 0, 0, 153, 154, 5, 47, 0, 0, 154, 155, 5, 47, 0, 0, 155,
		159, 1, 0, 0, 0, 156, 158, 9, 0, 0, 0, 157, 156, 1, 0, 0, 0, 158, 161,
		1, 0, 0, 0, 159, 160, 1, 0, 0, 0, 159, 157, 1, 0, 0, 0, 160, 163, 1, 0,
		0, 0, 161, 159, 1, 0, 0, 0, 162, 164, 5, 13, 0, 0, 163, 162, 1, 0, 0, 0,
		163, 164, 1, 0, 0, 0, 164, 165, 1, 0, 0, 0, 165, 166, 5, 10, 0, 0, 166,
		167, 1, 0, 0, 0, 167, 168, 6, 17, 0, 0, 168, 36, 1, 0, 0, 0, 169, 171,
		7, 5, 0, 0, 170, 169, 1, 0, 0, 0, 171, 172, 1, 0, 0, 0, 172, 170, 1, 0,
		0, 0, 172, 173, 1, 0, 0, 0, 173, 174, 1, 0, 0, 0, 174, 175, 6, 18, 0, 0,
		175, 38, 1, 0, 0, 0, 176, 178, 5, 13, 0, 0, 177, 176, 1, 0, 0, 0, 177,
		178, 1, 0, 0, 0, 178, 179, 1, 0, 0, 0, 179, 180, 5, 10, 0, 0, 180, 181,
		1, 0, 0, 0, 181, 182, 6, 19, 0, 0, 182, 40, 1, 0, 0, 0, 10, 0, 114, 120,
		126, 139, 145, 159, 163, 172, 177, 1, 6, 0, 0,
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
	TheoremLexerHEX       = 16
	TheoremLexerHEXDIG    = 17
	TheoremLexerCOMMENT   = 18
	TheoremLexerWS        = 19
	TheoremLexerNL        = 20
)
