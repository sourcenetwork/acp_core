package types

// FIXME refactor this, find a better error system

import (
	"errors"
	"fmt"
)

var ModuleName = "acp_core"

// x/acp module sentinel errors
var (
	// ErrAcpInternal is a general base error for IO or unexpected system errors
	ErrAcpInternal = errors.New("internal error")

	// ErrAcpInput is a general base error for input errors
	ErrAcpInput = errors.New("input error")

	// ErrAcpProtocolViolation is a general base error for operations forbidden by the protocol
	ErrAcpProtocolViolation = wrap("acp protocol violation", ErrAcpInput)

	// ErrAcpInvariantViolation indicates that an important condition of the protocol
	// has been violated, either by a bug or a successful exploit.
	// These are bad.
	ErrAcpInvariantViolation = errors.New("invariant violation")

	ErrPolicyNil        = wrap("policy must not be nil", ErrAcpInput)
	ErrRelationshipNil  = wrap("relationship must not be nil", ErrAcpInput)
	ErrActorNil         = wrap("actor must not be nil", ErrAcpInput)
	ErrRegistrationNil  = wrap("registration must not be nil", ErrAcpInput)
	ErrAccessRequestNil = wrap("AccessRequest must not be nil", ErrAcpInput)
	ErrInvalidVariant   = wrap("invalid type variant", ErrAcpInput)
	ErrObjectNil        = wrap("object must not be nil", ErrAcpInput)
	ErrTimestampNil     = wrap("timestamp must not be nil", ErrAcpInput)
	ErrAccNotFound      = wrap("account not found", ErrAcpInput)
	ErrPolicyNotFound   = wrap("policy not found", ErrAcpInput)
	ErrObjectNotFound   = wrap("object not found", ErrAcpInput)
	ErrInvalidHeight    = wrap("invalid block height", ErrAcpInput)
	ErrInvalidAccAddr   = wrap("invalid account address", ErrAcpInput)
	ErrInvalidDID       = wrap("invalid DID", ErrAcpInput)
	ErrAuthentication   = wrap("failed authentication check", ErrAcpInput)

	ErrNotAuthorized = wrap("actor not authorized", ErrAcpProtocolViolation)
)

func wrap(err string, base error) error {
	return fmt.Errorf("%v: %w", err, base)
}
