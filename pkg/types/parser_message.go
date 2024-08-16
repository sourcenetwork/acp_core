package types

import (
	fmt "fmt"

	"github.com/sourcenetwork/acp_core/pkg/errors"
)

func (e *LocatedMessage) ToError() error {
	return errors.New(fmt.Sprintf("%v:%v:%v %v",
		e.InputName, e.Range.Start.Line, e.Range.Start.Column, e.Message),
		errors.ErrorType_BAD_INPUT)
}

func (e *LocatedMessage) IsError() bool {
	return e.Kind == LocatedMessage_ERROR
}
