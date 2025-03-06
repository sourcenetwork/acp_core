package ppp

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

const defaultActorResourceName string = "actor"

var _ Transformer = (*BasicTransformer)(nil)

type BasicTransformer struct{}

func (s *BasicTransformer) Validate(pol *types.Policy) []error {
	var violations []error

	if pol.Name == "" {
		err := fmt.Errorf("name is required")
		violations = append(violations, err)
	}

	g := policy.BuildManagementGraph(pol)
	err := g.IsWellFormed()
	if err != nil {
		err := fmt.Errorf("%w: %w", policy.ErrInvalidManagementRule, err)
		violations = append(violations, err)
	}

	return violations
}

func (s *BasicTransformer) Transform(producer PolicyProvider) (*types.Policy, error) {
	pol := producer()

	if pol.ActorResource == nil {
		pol.ActorResource = &types.ActorResource{
			Name: defaultActorResourceName,
		}
	}

	// policy is sorted before building id to ensure determinism
	pol.Sort()

	return pol, nil
}

func (s *BasicTransformer) Name() string {
	return "Basic Transformer"
}
