package policy

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func HandleBuildCatalogue(ctx context.Context, manager runtime.RuntimeManager, polId string) (*types.Catalogue, error) {
	engine, err := zanzi.NewZanzi(manager.GetKVStore(), manager.GetLogger())
	if err != nil {
		return nil, err
	}

	rec, err := engine.GetPolicy(ctx, polId)
	if err != nil {
		return nil, err // TODO wrap
	}
	if rec == nil {
		return nil, errors.NewPolicyNotFound(polId)
	}

	catalogue := &types.Catalogue{
		PolicyDefinition:  rec.PolicyDefinition,
		ActorResourceName: rec.Policy.ActorResource.Name,
		ResourceCatalogue: make(map[string]*types.ResourceCatalogue),
	}

	for _, resource := range rec.Policy.Resources {
		resCatalogue := &types.ResourceCatalogue{}
		for _, permission := range resource.Permissions {
			resCatalogue.Permissions = append(resCatalogue.Permissions, permission.Name)
		}
		for _, relation := range resource.Relations {
			resCatalogue.Relations = append(resCatalogue.Relations, relation.Name)
		}
		catalogue.ResourceCatalogue[resource.Name] = resCatalogue
	}

	builder := types.RelationshipSelectorBuilder{}
	selector := builder.AnyObject().Relation(OwnerRelation).AnySubject().Build()
	ownerRelationships, err := engine.FilterRelationships(ctx, rec.Policy, &selector)
	if err != nil {
		return nil, err
	}

	for _, rec := range ownerRelationships {
		resCat := catalogue.ResourceCatalogue[rec.Relationship.Relation]
		resCat.ObjectIds = append(resCat.ObjectIds, rec.Relationship.Object.Id)
	}

	return catalogue, nil
}

/*

func buildCatalogue(ctx context.Context, data *data) (*types.Catalogue, error) {
	objs := coreutils.MapSlice(data.State.Objects, func(o *types.ObjectDefinition) *types.ObjectDetails {
		return &types.ObjectDetails{
			Resource:   o.Resource,
			Id:         o.Id,
			Owner:      data.ActorMap[o.Owner],
			OwnerAlias: o.Owner,
		}
	})

	aliases := make([]string, 0, len(data.State.Actors))
	actors := make([]*types.Actor, 0, len(data.State.Actors))
	for _, actorDef := range data.State.Actors {
		actor := &types.Actor{
			Id:    data.ActorMap[actorDef.Name],
			Alias: actorDef.Name,
		}
		actors = append(actors, actor)
		aliases = append(aliases, actorDef.Name)
	}
	actorset := &types.ActorSet{
		Aliases: aliases,
		Actors:  actors,
	}

	resCatalogues := make(map[string]*types.ResourceCatalogue)
	for _, resource := range data.Policy.Resources {
		catalogue, err := buildResourceCatalogue(ctx, data, resource, data.Engine)
		if err != nil {
			return nil, err
		}
		resCatalogues[resource.Name] = catalogue
	}

	relationships := coreutils.MapSlice(data.Relationships, func(relationships *acptypes.Relationship) string {
		return relationships.String() //TODO add pretty pritn method for relationship
	})

	catalogue := &types.Catalogue{
		Relationships:     relationships,
		Objects:           objs,
		Actors:            actorset,
		ResourceCatalogue: resCatalogues,
		ActorResourceName: data.Policy.ActorResource.Name,
		PolicyDefinition:  data.PolicyRaw,
	}

	return catalogue, nil
}

func buildResourceCatalogue(ctx context.Context, data *data, resource *acptypes.Resource, engine acptypes.ACPEngineServer) (*types.ResourceCatalogue, error) {
	builder := acptypes.RelationshipSelectorBuilder{}
	selector := builder.Resource(resource.Name).AnySubject().AnyRelation().Build()
	req := &acptypes.QueryFilterRelationshipsRequest{
		PolicyId: data.Policy.Id,
		Selector: &selector,
	}
	resp, err := engine.FilterRelationships(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("building catalogue: resource %v: fetching objects: %w", resource.Name, err)
	}

	ids := coreutils.MapSlice(resp.Records, func(rec *acptypes.RelationshipRecord) string { return rec.Relationship.Object.Id })

	return &types.ResourceCatalogue{
		Permissions: resource.ListPermissionsNames(),
		Relations:   resource.ListRelationsNames(),
		ObjectIds:   ids,
	}, nil
}

*/
