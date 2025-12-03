package ppp

import (
	"github.com/sourcenetwork/acp_core/internal/specification"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ specification.Requirement = (*BasicRequirement)(nil)

var ErrBasicRequirement = errors.New("basic requirement", errors.ErrorType_BAD_INPUT)

// BasicRequirement applies basic Policy Validation
// and performs basic Policy validation
type BasicRequirement struct{}

// Validate ensures that at minimum, the policy has a name
func (s *BasicRequirement) Validate(pol types.Policy) *errors.MultiError {
	if pol.Name == "" {
		return errors.NewMultiError(ErrBasicRequirement,
			errors.Wrap("name required", errors.ErrInvalidPolicy),
		)
	}

	return nil
}

func (t *BasicRequirement) GetName() string {
	return "Basic Requirement"
}
