package policy

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
)

var (
	ErrUnknownMarshalingType = errors.Wrap("unknown marshaling type", errors.ErrInput)
	ErrUnmarshaling          = errors.Wrap("unmarshaling error", errors.ErrInput)

	ErrInvalidShortPolicy           = errors.Wrap("invalid short policy", errors.ErrInput)
	ErrResourceMissingOwnerRelation = errors.Wrap("resource missing owner relation: %w", errors.ErrInvalidPolicy)
	ErrInvalidManagementRule        = errors.Wrap("invalid relation managamente definition: %w", errors.ErrInvalidPolicy)
)
