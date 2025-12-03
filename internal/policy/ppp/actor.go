package ppp

import (
	"github.com/sourcenetwork/acp_core/internal/specification"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ specification.Transformer = (*DefaultActorTransformer)(nil)

const ActorResourceName string = "actor"
const ActorResourceDoc = "actor resource models the set of actors defined within a policy"

var ErrActorResource = errors.New("actor resource not found", errors.ErrorType_BAD_INPUT)

// DefaultActorTransformer adds the default `actor` resource to a policy
type DefaultActorTransformer struct{}

// Validate ensures the Actor resource exists
func (s *DefaultActorTransformer) Validate(pol types.Policy) *errors.MultiError {
	if pol.ActorResource == nil {
		return errors.NewMultiError(ErrActorResource)
	}
	return nil
}

// Transform sets and creates the default ActorResource if ommitted
func (t *DefaultActorTransformer) Transform(pol types.Policy) (specification.TransformerResult, error) {
	result := specification.TransformerResult{}

	// pre-validate to ensure policy does not have any resource with clashing name
	for _, resource := range pol.Resources {
		if resource.Name == ActorResourceName {
			return result, errors.Wrap("policy defines reserved resource name `actor`: rename `actor` to a different name or remove the resource", errors.ErrInvalidPolicy)
		}
	}

	if pol.ActorResource == nil {
		pol.ActorResource = &types.ActorResource{}
	}
	pol.ActorResource.Name = ActorResourceName
	pol.ActorResource.Doc = ActorResourceDoc

	result.Policy = pol
	return result, nil
}

func (t *DefaultActorTransformer) GetName() string {
	return "Actor Transformer"
}
