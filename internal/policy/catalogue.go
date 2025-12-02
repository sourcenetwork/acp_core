package policy

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
	"k8s.io/apimachinery/pkg/util/sets"
)

// BuildCatalogue builds a new PolicyCatalogue from a PolicyID.
// Effectively it queries *all* relationships in a Policy and builds an index
// of known objets by resource.
// This is could be an expensive operation if the Policy has a lot of relationships.
func BuildCatalogue(ctx context.Context, engine *zanzi.Adapter, polId string) (*types.PolicyCatalogue, error) {
	// FIXME I should use indexes for this: index of registered users and objects
	rec, err := engine.GetPolicy(ctx, polId)
	if err != nil {
		return nil, errors.Wrap("fetching policy", err, errors.Pair("policy_id", polId))
	}
	if rec == nil {
		return nil, errors.ErrPolicyNotFound(polId)
	}

	actors := sets.New[string]()
	catalogue := &types.PolicyCatalogue{
		ActorResourceName: rec.Policy.ActorResource.Name,
		ResourceCatalogue: make(map[string]*types.ResourceCatalogue),
	}

	for _, resource := range rec.Policy.Resources {
		resCatalogue := &types.ResourceCatalogue{}

		resCatalogue.Permissions = utils.MapSlice(resource.Permissions, func(p *types.Permission) string { return p.Name })
		resCatalogue.Relations = utils.MapSlice(resource.GetAllRelations(), func(p *types.Relation) string { return p.Name })
		catalogue.ResourceCatalogue[resource.Name] = resCatalogue
	}

	for _, res := range rec.Policy.Resources {
		objects := sets.New[string]()

		builder := types.RelationshipSelectorBuilder{}
		selector := builder.Resource(res.Name).AnyRelation().AnySubject().Build()
		records, err := engine.FilterRelationships(ctx, rec.Policy, &selector)
		if err != nil {
			return nil, errors.Wrap("filtering relationships", err, errors.Pair("polId", polId))
		}
		for _, record := range records {
			if actor := record.Relationship.Subject.GetActor(); actor != nil {
				actors.Insert(actor.Id)
			}
			objects.Insert(record.Relationship.Object.Id)
		}

		catalogue.ResourceCatalogue[res.Name].ObjectIds = objects.UnsortedList()
	}

	catalogue.Actors = actors.UnsortedList()
	utils.SortSlice(catalogue.Actors)

	return catalogue, nil
}
