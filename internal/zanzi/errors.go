package zanzi

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
	zerrors "github.com/sourcenetwork/zanzi/pkg/errors"
)

func mapErr(err error) error {
	switch e := err.(type) {
	case *zerrors.Error:
		pairs := toAttrs(e.Metadata)
		return errors.NewWithCause(e.Message, e.Cause, mapKind(e.Kind), pairs...)
	case nil:
		return nil
	default:
		return errors.NewWithCause("zanzi error", err, errors.ErrorType_UNKNOWN)
	}
}

func mapKind(kind zerrors.ErrorKind) errors.ErrorType {
	switch kind {
	case zerrors.BadInput:
		return errors.ErrorType_BAD_INPUT
	case zerrors.Internal:
		return errors.ErrorType_INTERNAL
	case zerrors.NotFound:
		return errors.ErrorType_NOT_FOUND
	case zerrors.Unauthorized:
		return errors.ErrorType_UNAUTHORIZED
	case zerrors.Unknown:
		return errors.ErrorType_UNKNOWN
	case zerrors.Violation:
		return errors.ErrorType_OPERATION_FORBIDDEN
	default:
		return errors.ErrorType_UNKNOWN
	}
}

func toAttrs(attrs map[string]string) []errors.ContextPair {
	pairs := make([]errors.ContextPair, 0, len(attrs))
	for key, val := range attrs {
		pairs = append(pairs, errors.Pair(key, val))
	}
	return pairs
}
