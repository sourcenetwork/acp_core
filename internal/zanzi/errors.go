package zanzi

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
	zerrors "github.com/sourcenetwork/zanzi/pkg/errors"
)

func mapErr(err error) error {
	switch e := err.(type) {
	case *zerrors.Error:
		return errors.NewFromBaseError(e.Cause, mapKind(e.Kind), e.Message)
	case nil:
		return nil
	default:
		return errors.NewFromBaseError(err, errors.ErrorType_UNKNOWN, "zanzi errors")
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
