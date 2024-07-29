package errors

import (
	fmt "fmt"
	"strings"
)

var _ error = (*ParserMessage)(nil)
var _ TypedError = (*ParserMessage)(nil)
var _ error = (*ParserReport)(nil)
var _ TypedError = (*ParserReport)(nil)

func (e *ParserMessage) Error() string {
	return fmt.Sprintf("%v:%v:%v %v", e.InputName, e.Range.Start.Line, e.Range.Start.Column, e.Message)
}

func (e *ParserMessage) GetType() ErrorType {
	return ErrorType_BAD_INPUT
}

func (e *ParserReport) Error() string {
	builder := strings.Builder{}
	for _, err := range e.Messages {
		builder.WriteString(err.Error())
		builder.WriteRune('\n')
	}
	return builder.String()
}

func (e *ParserReport) GetType() ErrorType {
	return ErrorType_BAD_INPUT
}
