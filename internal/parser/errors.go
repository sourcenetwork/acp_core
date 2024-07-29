package parser

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/sourcenetwork/acp_core/pkg/errors"
)

var _ antlr.ErrorListener = (*errListener)(nil)

type errListener struct {
	errors errors.ParserReport
}

func (l *errListener) GetError() *errors.ParserReport {
	if l.errors.Messages == nil || len(l.errors.Messages) == 0 {
		return nil
	}

	return &l.errors
}

func (l *errListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	err := &errors.ParserMessage{
		Message: msg,
		Range: &errors.BufferRange{
			Start: &errors.BufferPosition{
				Line:   uint64(line),
				Column: uint64(column),
			},
			// Antlr doesn't provide the end position for the error,
			// default to 0,0 position
			End: &errors.BufferPosition{
				Line:   0,
				Column: 0,
			},
		},
	}
	l.errors.Messages = append(l.errors.Messages, err)

}

// TODO maybe?
func (l *errListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (l *errListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (l *errListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
}
