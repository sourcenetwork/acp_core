package ppp

import (
	"github.com/sourcenetwork/acp_core/internal/specification"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ specification.Transformer = (*BasicTransformer)(nil)
var _ specification.Requirement = (*BasicRequirement)(nil)

const ActorResourceName string = "actor"
const ActorResourceDoc = "actor resource models the set of actors defined within a policy"

// ErrBasicTransformer is the base error for problems detected by the BasicTransformer
var ErrBasicTransformer = errors.New("basic transformer", errors.ErrorType_BAD_INPUT)
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

// BasicTransforms normalizes a Policy by adding defaults
// to some optional fields
type BasicTransformer struct{}

// Validate ensures the Actor resource exists
func (s *BasicTransformer) Validate(pol types.Policy) *errors.MultiError {
	if pol.ActorResource == nil {
		return errors.NewMultiError(ErrBasicTransformer,
			errors.Wrap("invalid actor resource", errors.ErrInvalidPolicy),
		)
	}
	return nil
}

// Transform sets and creates the default ActorResource if ommitted
func (t *BasicTransformer) Transform(pol types.Policy) (specification.TransformerResult, error) {
	result := specification.TransformerResult{}
	if pol.ActorResource == nil {
		pol.ActorResource = &types.ActorResource{}
	}
	pol.ActorResource.Name = ActorResourceName
	pol.ActorResource.Doc = ActorResourceDoc

	result.Policy = pol
	return result, nil
}

func (t *BasicTransformer) GetName() string {
	return "Basic Transformer"
}
