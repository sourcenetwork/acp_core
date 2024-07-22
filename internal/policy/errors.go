package policy

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
)

var (
	ErrUnknownMarshalingType = errors.Wrap("unknown marshaling type", errors.ErrorType_BAD_INPUT)
	ErrUnmarshaling          = errors.Wrap("unmarshaling error", errors.ErrorType_BAD_INPUT)

	ErrInvalidShortPolicy           = errors.Wrap("invalid short policy", errors.ErrorType_BAD_INPUT)
	ErrResourceMissingOwnerRelation = errors.Wrap("resource missing owner relation: %w", errors.ErrInvalidPolicy)
	ErrInvalidManagementRule        = errors.Wrap("invalid relation managamente definition: %w", errors.ErrInvalidPolicy)
)
