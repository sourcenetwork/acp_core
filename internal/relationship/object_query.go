package relationship

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func GetObjectRegistrationHandler(ctx context.Context, runtime runtime.RuntimeManager, req *types.GetObjectRegistrationRequest) (*types.GetObjectRegistrationResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, newGetObjectRegistrationErr(err)
	}

	rec, err := engine.GetPolicy(ctx, req.PolicyId)
	if err != nil {
		return nil, newGetObjectRegistrationErr(err)
	}
	if rec == nil {
		return nil, newGetObjectRegistrationErr(errors.NewPolicyNotFound(req.PolicyId))
	}

	record, err := queryOwnerRelationship(ctx, engine, rec.Policy, req.Object)
	if err != nil {
		return nil, newGetObjectRegistrationErr(err)
	}
	if record == nil {
		return &types.GetObjectRegistrationResponse{
			IsRegistered: false,
			OwnerId:      "",
		}, nil
	}

	// Currently only Actors should be Object owners,
	// therefore if an `owner` relationship was found it must be an actor.
	// Nevertheless, in the off chance the owner isn't an Actor,
	// Return an error and log it.
	actor := record.Relationship.Subject.GetActor()
	if actor == nil {
		// TODO Emit metric and Log
		return nil, newGetObjectRegistrationErr(errors.Wrap("object owner isn't an actor", errors.ErrInvariantViolation,
			errors.Pair("policy", req.PolicyId),
			errors.Pair("resource", req.Object.Resource),
			errors.Pair("object", req.Object.Id),
		))
	}
	return &types.GetObjectRegistrationResponse{
		IsRegistered: true,
		OwnerId:      actor.Id,
		Record:       record,
	}, nil
}
