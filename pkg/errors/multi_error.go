package errors

import "strings"

var _ error = (*MultiError)(nil)

// MultiError models an aggregate of errors
// with a common underlying error cause
type MultiError struct {
	msg  string
	errs []error
	kind ErrorType
}

// GetErrors return the individual errors in MultiError
func (m *MultiError) GetErrors() []error {
	return m.errs
}

// GetType returns the MultiError ErrorType
func (m *MultiError) GetType() ErrorType {
	return m.kind
}

// Append adds errs to MultiErr
func (m *MultiError) Append(errs ...error) {
	m.errs = append(m.errs, errs...)
}

// Error implements Go's error interface
func (m *MultiError) Error() string {
	builder := strings.Builder{}
	builder.WriteString(m.msg)
	builder.WriteRune(':')
	for _, err := range m.errs {
		builder.WriteString(err.Error())
		builder.WriteString("; ")
	}
	return builder.String()
}

func NewMultiError(msg string, kind ErrorType, errs ...error) *MultiError {
	return &MultiError{
		errs: errs,
		kind: kind,
		msg:  msg,
	}
}
