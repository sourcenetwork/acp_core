// Code generated from ./pkg/parser/permission_parser/PermissionExpr.g4 by ANTLR 4.13.2. DO NOT EDIT.

package permission_parser

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

type PermissionExprLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var PermissionExprLexerLexerStaticData struct {
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

func permissionexprlexerLexerInit() {
	staticData := &PermissionExprLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.LiteralNames = []string{
		"", "'->'", "'('", "')'", "'+'", "'-'", "'&'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "", "", "", "IDENTIFIER", "WS",
	}
	staticData.RuleNames = []string{
		"T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "IDENTIFIER", "WS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 8, 44, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 1, 0, 1, 0, 1, 0, 1, 1, 1,
		1, 1, 2, 1, 2, 1, 3, 1, 3, 1, 4, 1, 4, 1, 5, 1, 5, 1, 6, 1, 6, 5, 6, 33,
		8, 6, 10, 6, 12, 6, 36, 9, 6, 1, 7, 4, 7, 39, 8, 7, 11, 7, 12, 7, 40, 1,
		7, 1, 7, 0, 0, 8, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 1,
		0, 3, 2, 0, 65, 90, 97, 122, 4, 0, 48, 57, 65, 90, 95, 95, 97, 122, 2,
		0, 9, 9, 32, 32, 45, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0,
		0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0,
		0, 0, 0, 15, 1, 0, 0, 0, 1, 17, 1, 0, 0, 0, 3, 20, 1, 0, 0, 0, 5, 22, 1,
		0, 0, 0, 7, 24, 1, 0, 0, 0, 9, 26, 1, 0, 0, 0, 11, 28, 1, 0, 0, 0, 13,
		30, 1, 0, 0, 0, 15, 38, 1, 0, 0, 0, 17, 18, 5, 45, 0, 0, 18, 19, 5, 62,
		0, 0, 19, 2, 1, 0, 0, 0, 20, 21, 5, 40, 0, 0, 21, 4, 1, 0, 0, 0, 22, 23,
		5, 41, 0, 0, 23, 6, 1, 0, 0, 0, 24, 25, 5, 43, 0, 0, 25, 8, 1, 0, 0, 0,
		26, 27, 5, 45, 0, 0, 27, 10, 1, 0, 0, 0, 28, 29, 5, 38, 0, 0, 29, 12, 1,
		0, 0, 0, 30, 34, 7, 0, 0, 0, 31, 33, 7, 1, 0, 0, 32, 31, 1, 0, 0, 0, 33,
		36, 1, 0, 0, 0, 34, 32, 1, 0, 0, 0, 34, 35, 1, 0, 0, 0, 35, 14, 1, 0, 0,
		0, 36, 34, 1, 0, 0, 0, 37, 39, 7, 2, 0, 0, 38, 37, 1, 0, 0, 0, 39, 40,
		1, 0, 0, 0, 40, 38, 1, 0, 0, 0, 40, 41, 1, 0, 0, 0, 41, 42, 1, 0, 0, 0,
		42, 43, 6, 7, 0, 0, 43, 16, 1, 0, 0, 0, 3, 0, 34, 40, 1, 6, 0, 0,
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

// PermissionExprLexerInit initializes any static state used to implement PermissionExprLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewPermissionExprLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func PermissionExprLexerInit() {
	staticData := &PermissionExprLexerLexerStaticData
	staticData.once.Do(permissionexprlexerLexerInit)
}

// NewPermissionExprLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewPermissionExprLexer(input antlr.CharStream) *PermissionExprLexer {
	PermissionExprLexerInit()
	l := new(PermissionExprLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &PermissionExprLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "PermissionExpr.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// PermissionExprLexer tokens.
const (
	PermissionExprLexerT__0       = 1
	PermissionExprLexerT__1       = 2
	PermissionExprLexerT__2       = 3
	PermissionExprLexerT__3       = 4
	PermissionExprLexerT__4       = 5
	PermissionExprLexerT__5       = 6
	PermissionExprLexerIDENTIFIER = 7
	PermissionExprLexerWS         = 8
)
