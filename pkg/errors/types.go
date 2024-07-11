package errors

import "fmt"

// ContextPair is a key value pair used to annotate the context under which the error was found
type ContextPair struct {
	key string
	val any
}

// Pair returns a new ContextPair
func Pair(key string, val any) ContextPair {
	return ContextPair{
		key: key,
		val: val,
	}
}

type ErrorType string

const (
	// ErrUnspecified is a catch all bucket for unspecified errors
	ErrUnspecified ErrorType = ErrorType("uknown internal error")
	// ErrInternal is a general base error for IO or unexpected system errors
	ErrInternal = ErrorType("internal error")
	// ErrUnauthenticated signals caller was not authenticated while it was required
	ErrUnauthenticated = ErrorType("not authenticated")
	// ErrUnauthorized is a general error for operations that were
	// rejected due to insufficient permission unauthorized
	ErrUnauthorized = ErrorType("unauthorized")
	// Errinput is a general base error for input errors
	ErrInput = ErrorType("input error")
	// ErrOperationForbidden signals that the operation was not executed
	// as that would violate part of the Access Control system
	ErrOperationForbidden = ErrorType("forbidden")
)

func (t ErrorType) String() string {
	return string(t)
}

func (t ErrorType) Error() string {
	return t.String()
}

func (t ErrorType) Is(err error) bool {
	return err.Error() == t.Error()
}

type Error struct {
	errType ErrorType
	base    error
	message string
	pairs   []ContextPair
}

func (e *Error) Unwrap() []error {
	return []error{e.base, e.errType}
}

func (e *Error) getMsgChain() string {
	if e.base == nil {
		return e.message
	}
	switch err := e.base.(type) {
	case nil:
		return e.message
	case *Error:
		return fmt.Sprintf("%v: %v", e.message, err.getMsgChain())
	default:
		return fmt.Sprintf("%v: %v", e.message, e.base.Error())
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v: ctx={%v}", e.getMsgChain(), e.pairs)
}

// AppendPairs returns a new Error with the extra information contained in pairs
// Sets the current error as the base error for the new error
func (e *Error) AppendPairs(pairs ...ContextPair) *Error {
	return &Error{
		errType: e.errType,
		base:    e,
		message: e.message,
		pairs:   append(e.pairs, pairs...),
	}
}

// Refine returns a new error with the additional context data.
// Preserves the underlying base error
func (e *Error) Refine(message string, pairs ...ContextPair) *Error {
	return &Error{
		errType: e.errType,
		base:    e,
		message: message,
		pairs:   append(e.pairs, pairs...),
	}
}

// NewFromCause creates a new error with the aditional context given and formats the message to include
// the cause error. It does not add cause to the error chain.
func NewFromCause(msg string, cause error, errType ErrorType, pairs ...ContextPair) *Error {
	return &Error{
		errType: errType,
		message: fmt.Sprintf("%v: %v", msg, cause),
		pairs:   pairs,
	}
}

// NewFromBaseError returns a new Error which wraps the given error, including
// the additional context data given.
// Includes base as part of the error chain
func NewFromBaseError(base error, errType ErrorType, msg string, pairs ...ContextPair) *Error {
	return &Error{
		errType: errType,
		base:    base,
		message: msg,
		pairs:   pairs,
	}
}

// New creates a new error from a message, an ErrorType and an optional set of context pairs
func New(message string, errType ErrorType, pairs ...ContextPair) *Error {
	return &Error{
		errType: errType,
		base:    errType,
		message: message,
		pairs:   pairs,
	}
}

// Wrap refines a given error with the additional message and context mesages,
// includes err as part of the error chain.
func Wrap(message string, err error, pairs ...ContextPair) error {
	switch e := err.(type) {
	case *Error:
		return e.Refine(message, pairs...)
	case ErrorType:
		return New(message, e)
	default:
		return NewFromBaseError(err, ErrUnspecified, message, pairs...)
	}
}
