package relationship

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

var (
	ErrDeleteOwnerRel      = fmt.Errorf("cannot delete an owner relationship: %w", types.ErrAcpProtocolViolation)
	ErrSetOwnerRel         = fmt.Errorf("cannot set an owner relationship: %w", types.ErrAcpProtocolViolation)
	ErrInvalidRelationship = fmt.Errorf("invalid relationship: %w", types.ErrAcpInput)
)
