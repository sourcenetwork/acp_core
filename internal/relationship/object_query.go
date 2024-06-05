package relationship

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func GetObjectRegistrationHandler(ctx context.Context, runtime runtime.RuntimeManager, req *types.GetObjectRegistrationRequest) (*types.GetObjectRegistrationResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	rec, err := engine.GetPolicy(ctx, req.PolicyId)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, fmt.Errorf("policy %v: %w", req.PolicyId, types.ErrPolicyNotFound)
	}

	builder := types.RelationshipSelectorBuilder{}
	builder.Object(req.Object)
	builder.Relation(policy.OwnerRelation)
	builder.AnySubject()
	selector := builder.Build()

	records, err := engine.FilterRelationships(ctx, rec.Policy, &selector)
	if err != nil {
		return nil, err
	}

	response := &types.GetObjectRegistrationResponse{}

	if len(records) > 0 {
		// Currently only Actors should be Object owners,
		// therefore if an `owner` relationship was found it must be an actor.
		// Nevertheless, in the off chance the owner isn't an Actor,
		// Return an error and log it.
		actor := records[0].Relationship.Subject.GetActor()
		if actor == nil {
			// TODO Emit metric
			runtime.GetLogger().Error("invariant error: object owner isn't type actor", "policyId", req.PolicyId, "object", req.Object, "relationship", records[0])
			return nil, fmt.Errorf("object %v has non Actor type as owner: %v", req.Object.Id, types.ErrAcpProtocolViolation)
		}

		response.OwnerId = actor.Id
		response.IsRegistered = true
	}

	return response, nil

}
