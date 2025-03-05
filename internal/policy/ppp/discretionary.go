package ppp

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

const OwnerRelationName = "owner"
const OwnerDescription = "owner relations represents the object owner"

var _ Transformer = (*DiscretionaryTransformer)(nil)

type DiscretionaryTransformer struct{}

func (t *DiscretionaryTransformer) Name() string {
	return "Discretionary Transformer"
}

func (t *DiscretionaryTransformer) Validate(policy *types.Policy) []error {
	var violations []error
	for _, resource := range policy.Resources {
		ownerRel := utils.FilterSlice(resource.Relations, func(r *types.Relation) bool { return r.Name == OwnerRelationName })
		if len(ownerRel) > 1 {
			err := fmt.Errorf("invalid policy: resource %v: multiple owner relations", resource.Name)
			violations = append(violations, err)
		}

		if len(ownerRel) == 0 {
			err := fmt.Errorf("invalid policy: resource %v: no owner relation")
			violations = append(violations, err)
		}
	}

	// TODO validate the top level node of all permissiosn is the owner

	return violations
}

func (t *DiscretionaryTransformer) Transform(provider PolicyProvider) (*types.Policy, error) {
	policy := provider()

	// add owner as top level relation in permission tree
	// add owner relation

	for _, resource := range policy.Resources {
		ownerRel := utils.FilterSlice(resource.Relations, func(r *types.Relation) bool { return r.Name == OwnerRelationName })
		if len(ownerRel) > 1 {
			return nil, fmt.Errorf("invalid policy: resource %v: multiple owner relations", resource.Name)
		}
		if len(ownerRel) == 0 {
			rel := newOwnerRelation(policy)
			ownerRel = append(ownerRel, rel)
			resource.Relations = append(resource.Relations, rel)
		}
	}

	for _, resource := range policy.Resources {
		for _, permission := range resource.Permissions {
			// FIXME this is problematic, i only have the string expression. frig
			_ = permission
		}
	}

	return policy, nil

}

func newOwnerRelation(policy *types.Policy) *types.Relation {
	return &types.Relation{
		Name: OwnerRelationName,
		Doc:  OwnerDescription,
		VrTypes: []*types.Restriction{
			{
				ResourceName: policy.ActorResource.Name,
			},
		},
	}
}
