package ppp

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/specification"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ specification.Transformer = (*BasicTransformer)(nil)
var _ specification.Requirement = (*BasicRequirement)(nil)

const DefaultActorResourceName string = "actor"

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
	if pol.ActorResource == nil || pol.ActorResource.Name == "" {
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
		pol.ActorResource = &types.ActorResource{
			Name: DefaultActorResourceName,
		}
	}

	// normalize all empty permissions to "owner"
	for _, res := range pol.Resources {
		for _, perm := range res.Permissions {
			if perm.Expression == "" {
				perm.Expression = "owner"
				msg := fmt.Sprintf("defaulting permission to owner: resource %v: relation %v",
					res.Name, perm.Name)
				result.Messages = append(result.Messages, msg)
			}
		}
	}
	result.Policy = pol
	return result, nil
}

func (t *BasicTransformer) GetName() string {
	return "Basic Transformer"
}
