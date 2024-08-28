//go:build js

package errors

import "syscall/js"

// NewJSError receives a Go Error and transforms it into a JS universe Error.
// The JS Error message is the same as the Go's Error() and its name
// matches err's ErrorType or ErrorType_UNKNOWN if not set.
func NewJSError(err error) js.Value {
	name := ErrorType_UNKNOWN.String()
	if e, ok := err.(TypedError); ok {
		name = e.GetType().String()
	}

	jsErr := js.Global().Get("Error").New(err.Error())
	jsErr.Set("name", name)
	jsErr.Set("message", err.Error())
	return jsErr
}
