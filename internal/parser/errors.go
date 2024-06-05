package parser

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

type ParseErrors struct {
	errors []ErrorDetail
}

func (p *ParseErrors) GetErrors() []ErrorDetail {
	return p.errors
}

func (p *ParseErrors) Error() string {
	builder := strings.Builder{}
	for _, err := range p.errors {
		builder.WriteString(err.Error())
		builder.WriteRune('\n')
	}

	return builder.String()
}

type ErrorDetail struct {
	Line   uint
	Column uint
	Msg    string
}

func (e *ErrorDetail) Error() string {
	return fmt.Sprintf("line %v:%v %v", e.Line, e.Column, e.Msg)
}

var _ antlr.ErrorListener = (*errListener)(nil)

type errListener struct {
	errors []ErrorDetail
}

func (l *errListener) GetError() *ParseErrors {
	if l.errors == nil || len(l.errors) == 0 {
		return nil
	}

	return &ParseErrors{
		errors: l.errors,
	}
}

func (l *errListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	err := ErrorDetail{
		Line:   uint(line),
		Column: uint(column),
		Msg:    msg,
	}
	l.errors = append(l.errors, err)

}

func (l *errListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}
func (l *errListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
}
func (l *errListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
}
