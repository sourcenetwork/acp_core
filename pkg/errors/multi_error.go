package errors

import "strings"

var _ error = (*MultiError)(nil)

// MultiError models an aggregate of errors
// with a common underlying error cause
type MultiError struct {
	root *Error
	errs []error
}

// GetErrors return the individual errors in MultiError
func (m *MultiError) GetErrors() []error {
	return m.errs
}

// GetRoot returns the MultiError's root *Error
func (m *MultiError) GetRoot() *Error {
	return m.root
}

// Append adds errs to MultiErr
func (m *MultiError) Append(errs ...error) {
	m.errs = append(m.errs, errs...)
}

// Error implements Go's error interface
func (m *MultiError) Error() string {
	builder := strings.Builder{}
	builder.WriteString(m.root.Error())
	builder.WriteString(": ")
	for _, err := range m.errs {
		builder.WriteString(err.Error())
		builder.WriteString("; ")
	}
	return builder.String()
}

// Unwrap implements go implicit Unwrap interface
func (e *MultiError) Unwrap() []error {
	return append(e.errs, e.root)
}

func NewMultiError(root *Error, errs ...error) *MultiError {
	return &MultiError{
		root: root,
		errs: errs,
	}
}
