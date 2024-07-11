package relationship

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func queryOwnerRelationship(ctx context.Context, engine *zanzi.Adapter, pol *types.Policy, obj *types.Object) (*types.RelationshipRecord, error) {
	builder := types.RelationshipSelectorBuilder{}
	builder.Object(obj)
	builder.Relation(policy.OwnerRelation)
	builder.AnySubject()
	selector := builder.Build()

	records, err := engine.FilterRelationships(ctx, pol, &selector)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, nil
	} else if len(records) == 1 {
		return records[0], nil
	} else {
		// Every object should only have one owner.
		// In the event this anomaly happen, return the first owner but log the incident
		// and fire a new invariant violation
		// TODO Log and emit metricEmit metric
		//ctx.Logger().Error("invariant error: object owner isn't type actor", "policyId", req.PolicyId, "object", req.Object, "relationship", records[0])
		return records[0], nil
	}
}
