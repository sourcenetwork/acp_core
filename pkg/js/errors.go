//go:build js

package js

import (
	"syscall/js"

	"github.com/sourcenetwork/acp_core/pkg/errors"
)

func newInvalidArgsErr(count int) error {
	return errors.Wrap("invalid number of arguments", errors.ErrorType_BAD_INPUT, errors.Pair("count", count))
}

func newJSError(err error) js.Value {
	name := errors.ErrorType_UNKNOWN.String()
	if e, ok := err.(errors.TypedError); ok {
		name = e.GetType().String()
	}

	jsErr := js.Global().Get("Error").New(err.Error())
	jsErr.Set("name", name)
	jsErr.Set("message", err.Error())
	return jsErr
}
