package errors

// Wrap attenuates an error type by adding an extra error message
// as well as a set of key-value metadata pairs
func Wrap(msg string, err error, pairs ...ContextPair) error {
	switch e := err.(type) {
	case ErrorType:
		return &Error{
			Kind:     e,
			Message:  msg,
			Cause:    nil,
			Metadata: appendPairs(nil, pairs),
		}
	case *Error:
		return &Error{
			Kind:     e.Kind,
			Message:  msg + ": " + e.Message,
			Cause:    e.Cause,
			Metadata: appendPairs(e.Metadata, pairs),
		}
	default:
		return &Error{
			Kind:     ErrorType_UNKNOWN,
			Message:  msg,
			Cause:    err,
			Metadata: nil,
		}
	}
}

// Attrs attenuates an error by adding additional key-value context metadata.
// Overwrites keys with the same name
func Attrs(err error, pairs ...ContextPair) error {
	switch e := err.(type) {
	case ErrorType:
		return &Error{
			Kind:     e,
			Message:  "",
			Cause:    nil,
			Metadata: appendPairs(nil, pairs),
		}
	case *Error:
		return &Error{
			Kind:     e.Kind,
			Message:  e.Message,
			Cause:    e.Cause,
			Metadata: appendPairs(e.Metadata, pairs),
		}
	default:
		return &Error{
			Kind:     ErrorType_UNKNOWN,
			Message:  "",
			Cause:    err,
			Metadata: appendPairs(nil, pairs),
		}
	}
}
