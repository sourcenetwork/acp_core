package errors

var _ error = ErrorType(ErrorType_UNKNOWN)

func (e ErrorType) Error() string {
	return e.String()
}

func (e ErrorType) Code() int {
	return int(e)
}

func (e ErrorType) Is(err error) bool {
	if castErr, ok := err.(ErrorType); ok {
		return castErr == e
	}
	return false
}
