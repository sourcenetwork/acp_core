package policy

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

// BuildCatalogue builds a new PolicyCatalogue from a PolicyID.
// Effectively it queries *all* registration relationships in a Policy and builds an index
// of known objets by resource.
// This is could be an expensive operation if the Policy has a lot of relationships.
func BuildCatalogue(ctx context.Context, engine *zanzi.Adapter, polId string) (*types.PolicyCatalogue, error) {
	rec, err := engine.GetPolicy(ctx, polId)
	if err != nil {
		return nil, errors.Wrap("fetching policy", err, errors.Pair("policy_id", polId))
	}
	if rec == nil {
		return nil, errors.ErrPolicyNotFound(polId)
	}

	actorSet := make(map[string]struct{})
	catalogue := &types.PolicyCatalogue{
		ActorResourceName: rec.Policy.ActorResource.Name,
		ResourceCatalogue: make(map[string]*types.ResourceCatalogue),
	}

	for _, resource := range rec.Policy.Resources {
		resCatalogue := &types.ResourceCatalogue{}
		resCatalogue.Permissions = utils.MapSlice(resource.Permissions, func(p *types.Permission) string { return p.Name })
		resCatalogue.Relations = utils.MapSlice(resource.Relations, func(p *types.Relation) string { return p.Name })
		catalogue.ResourceCatalogue[resource.Name] = resCatalogue
	}

	builder := types.RelationshipSelectorBuilder{}
	selector := builder.AnyObject().Relation(OwnerRelation).AnySubject().Build()
	ownerRelationships, err := engine.FilterRelationships(ctx, rec.Policy, &selector)
	if err != nil {
		return nil, errors.Wrap("filtering owner relationships", err, errors.Pair("polId", polId))
	}

	for _, rec := range ownerRelationships {
		resCat := catalogue.ResourceCatalogue[rec.Relationship.Object.Resource]
		resCat.ObjectIds = append(resCat.ObjectIds, rec.Relationship.Object.Id)
		actorSet[rec.Metadata.Creator.Identifier] = struct{}{}
	}

	for actor := range actorSet {
		catalogue.Actors = append(catalogue.Actors, actor)
	}
	utils.SortSlice(catalogue.Actors)

	return catalogue, nil
}
