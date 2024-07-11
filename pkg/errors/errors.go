package errors

var (
	// ErrInvariantViolation indicates that an important condition of the protocol
	// has been violated, either by a bug or a successful exploit.
	ErrInvariantViolation = New("invariant violation", ErrInternal)

	// ErrAcpProtocolViolation is a general base error for operations forbidden by the protocol
	//ErrAcpProtocolViolation = wrap("acp protocol violation", Errinput)

	// ErrNotFound signals some system object was not found
	ErrNotFound = New("not found", ErrInput)

	// ErrInvalidData signals some user supplied data was invalid
	ErrInvalidData = New("invalid data", ErrInput)

	// ErrUnknownVariant signals that an enum-like field received an unexpected value
	ErrUnknownVariant = New("unknown variant", ErrInput)

	ErrInvalidDID          = New("did", ErrInput)
	ErrInvalidPolicy       = New("policy", ErrInput)
	ErrInvalidRelationship = New("relationship", ErrInput)
)

func NewPolicyNotFound(policyID string) error {
	return Wrap("policy not found", ErrNotFound, Pair("policy", policyID))
}

func NewObjectNotFound(policyID, resource, object string) error {
	return Wrap("object not found", ErrNotFound,
		Pair("policy", policyID),
		Pair("resource", resource),
		Pair("objId", object),
	)
}
