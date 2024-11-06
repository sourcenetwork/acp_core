package parser

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ antlr.ErrorListener = (*errListener)(nil)

type errListener struct {
	report *ParserReport
}

func newListener(productionRuleName string) *errListener {
	return &errListener{
		report: &ParserReport{
			msg:  fmt.Sprintf("%v parser report", productionRuleName),
			msgs: nil,
		},
	}
}

// GetERror produces a ParserReport from the errors stored in the listener
func (l *errListener) GetReport() *ParserReport {
	return l.report
}

func (l *errListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	err := &types.LocatedMessage{
		Message: msg,
		Kind:    types.LocatedMessage_ERROR,
		Interval: &types.BufferInterval{
			Start: &types.BufferPosition{
				Line:   uint64(line),
				Column: uint64(column),
			},
			// Antlr doesn't provide the end position for the error,
			// default to 0,0 position
			// which can understood as EOF
			End: &types.BufferPosition{
				Line:   0,
				Column: 0,
			},
		},
	}
	l.report.msgs = append(l.report.msgs, err)
}

func (l *errListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (l *errListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}

func (l *errListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
}
