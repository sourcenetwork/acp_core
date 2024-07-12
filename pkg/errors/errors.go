package errors

var (
	// ErrInvariantViolation indicates that an important condition of the protocol
	// has been violated, either by a bug or a successful exploit.
	ErrInvariantViolation = New("invariant violation", ErrorType_INTERNAL)

	// ErrInvalidData signals some user supplied data was invalid
	ErrInvalidData = New("invalid data", ErrorType_BAD_INPUT)

	// ErrUnknownVariant signals that an enum-like field received an unexpected value
	ErrUnknownVariant = New("unknown variant", ErrorType_BAD_INPUT)

	ErrInvalidDID          = New("did", ErrorType_BAD_INPUT)
	ErrInvalidPolicy       = New("policy", ErrorType_BAD_INPUT)
	ErrInvalidRelationship = New("relationship", ErrorType_BAD_INPUT)
)

func NewPolicyNotFound(policyID string) error {
	return New("policy not found", ErrorType_NOT_FOUND, Pair("policy", policyID))
}

func NewObjectNotFound(policyID, resource, object string) error {
	return New("object not found", ErrorType_NOT_FOUND,
		Pair("policy", policyID),
		Pair("resource", resource),
		Pair("objId", object),
	)
}
