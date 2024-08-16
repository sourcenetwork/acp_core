package parser

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ antlr.ErrorListener = (*errListener)(nil)

type errListener struct {
	msgs []*types.LocatedMessage
}

// GetERror produces a ParserReport from the errors stored in the listener
func (l *errListener) GetMessages() []*types.LocatedMessage {
	return l.msgs
}

func (l *errListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	err := &types.LocatedMessage{
		Message: msg,
		Range: &types.BufferRange{
			Start: &types.BufferPosition{
				Line:   uint64(line),
				Column: uint64(column),
			},
			// Antlr doesn't provide the end position for the error,
			// default to 0,0 position
			End: &types.BufferPosition{
				Line:   0,
				Column: 0,
			},
		},
	}
	l.msgs = append(l.msgs, err)
}

func (l *errListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (l *errListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (l *errListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
}
