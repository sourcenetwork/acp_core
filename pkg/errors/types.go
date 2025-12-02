package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ContextPair model a key value attribute pair
type ContextPair struct {
	Key   string
	Value string
}

// Pair returns a new AttrPair with the given key value
func Pair(key, value string) ContextPair {
	return ContextPair{
		Key:   key,
		Value: value,
	}
}

var _ error = (*Error)(nil)

// NewWithCause returns a new zanzi error which wraps an ErrorType and an underlying cause error
func NewWithCause(msg string, cause error, kind ErrorType, pairs ...ContextPair) error {
	return &Error{
		Kind:     kind,
		Message:  msg,
		Cause:    cause,
		Metadata: appendPairs(nil, pairs),
	}
}

// New returns a new Zanzi base error type with the given kind
func New(message string, kind ErrorType, pairs ...ContextPair) *Error {
	return &Error{
		Kind:     kind,
		Message:  message,
		Cause:    nil,
		Metadata: appendPairs(nil, pairs),
	}
}

// Error models a Zanzi error object, which may wrap an underlaying cause error
// and contains a set of string key-value pairs which contain request specific metadata
// which caused the error
type Error struct {
	Kind     ErrorType
	Message  string
	Cause    error
	Metadata map[string]string
}

func (e *Error) Error() string {
	metadata, _ := json.Marshal(e.Metadata)
	str := fmt.Sprintf("%v; attrs=%v; kind=%v", e.Message, string(metadata), e.Kind.Error())
	if e.Cause != nil {
		str += "; cause: " + e.Cause.Error()
	}
	return str
}

// Is return true if target is of type ErrorType and e.Kind is the same as ErrorType
// or if target is also an Error instance and e's message contains target's message.
func (e *Error) Is(target error) bool {
	switch other := target.(type) {
	case *Error:
		return strings.Contains(e.Message, other.Message) && e.Kind == other.Kind || errors.Is(e.Cause, target)
	case ErrorType:
		return e.Kind.Is(other)
	default:
		return errors.Is(e.Cause, target) //not sure this is right
	}
}

func appendPairs(m map[string]string, pairs []ContextPair) map[string]string {
	attrs := make(map[string]string, len(m)+len(pairs))
	for k, v := range m {
		attrs[k] = v
	}
	for _, p := range pairs {
		attrs[p.Key] = p.Value
	}
	return attrs
}
