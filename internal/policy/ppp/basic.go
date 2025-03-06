package ppp

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

const defaultActorResourceName string = "actor"

var ErrInvalidManagementRule = errors.Wrap("invalid relation managament definition: %w", errors.ErrInvalidPolicy)
var _ Transformer = (*BasicTransformer)(nil)

type BasicTransformer struct{}

func (s *BasicTransformer) Validate(pol *types.Policy) *errors.MultiError {
	multiErr := errors.NewMultiError(ErrBasicTransformer)

	if pol.Name == "" {
		multiErr.Append(errors.Wrap("name required", errors.ErrInvalidPolicy))
	}

	g := &types.ManagementGraph{}
	g.LoadFromPolicy(pol)
	err := g.IsWellFormed()
	if err != nil {
		err := fmt.Errorf("%w: %w", ErrInvalidManagementRule, err)
		multiErr.Append(err)
	}

	if len(multiErr.GetErrors()) > 0 {
		return multiErr
	}
	return nil
}

// normalize normalizes a policy by setting default values for optional fields.
func (s *BasicTransformer) Transform(producer PolicyProvider) (*types.Policy, *errors.MultiError) {
	pol := producer()

	if pol.ActorResource == nil {
		pol.ActorResource = &types.ActorResource{
			Name: defaultActorResourceName,
		}
	}

	return pol, nil
}
