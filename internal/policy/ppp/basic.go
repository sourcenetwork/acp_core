package ppp

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/transformer"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ transformer.Transformer = (*BasicTransformer)(nil)
var _ transformer.Specification = (*BasicSpec)(nil)

const DefaultActorResourceName string = "actor"

// ErrBasicTransformer is the base error for problems detected by the BasicTransformer
var ErrBasicTransformer = errors.New("basic transformer", errors.ErrorType_BAD_INPUT)
var ErrBasicSpec = errors.New("basic spec", errors.ErrorType_BAD_INPUT)

// BasicSpec applies basic Policy Validation
// and performs basic Policy validation
type BasicSpec struct{}

// Validate ensures that at minimum, the policy has a name
func (s *BasicSpec) Validate(pol types.Policy) *errors.MultiError {
	if pol.Name == "" {
		return errors.NewMultiError(ErrBasicSpec,
			errors.Wrap("name required", errors.ErrInvalidPolicy),
		)
	}

	return nil
}

func (t *BasicSpec) GetBaseError() error {
	return ErrBasicSpec
}

// BasicTransforms normalizes a Policy by adding defaults
// to some optional fields
type BasicTransformer struct{}

// Validate ensures the Actor resource exists
func (s *BasicTransformer) Validate(pol types.Policy) *errors.MultiError {
	if pol.ActorResource == nil || pol.ActorResource.Name == "" {
		return errors.NewMultiError(ErrBasicTransformer,
			errors.Wrap("invalid actor resource", errors.ErrInvalidPolicy),
		)
	}
	return nil
}

// Transform sets and creates the default ActorResource if ommitted
func (t *BasicTransformer) Transform(pol types.Policy) (types.Policy, error) {
	if pol.ActorResource == nil {
		pol.ActorResource = &types.ActorResource{
			Name: DefaultActorResourceName,
		}
	}
	return pol, nil
}

func (t *BasicTransformer) GetBaseError() error {
	return ErrBasicTransformer
}
