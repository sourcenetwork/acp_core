package policy

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

var (
	ErrInvalidPolicy = fmt.Errorf("invalid policy: %w", types.ErrAcpInput)

	ErrUnknownMarshalingType = fmt.Errorf("unknown marshaling type: %w", types.ErrAcpInput)
	ErrUnmarshaling          = fmt.Errorf("unmarshaling error: %w", types.ErrAcpInput)

	ErrInvalidShortPolicy           = fmt.Errorf("invalid short policy: %w", ErrInvalidPolicy)
	ErrInvalidCreator               = fmt.Errorf("invalid creator: %w", ErrInvalidPolicy)
	ErrResourceMissingOwnerRelation = fmt.Errorf("resource missing owner relation: %w", ErrInvalidPolicy)
	ErrInvalidManagementRule        = fmt.Errorf("invalid relation managamente definition: %w", ErrInvalidPolicy)
)
