package simulator

import (
	"context"

	prototypes "github.com/cosmos/gogoproto/types"

	"github.com/sourcenetwork/acp_core/internal/parser"
	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/internal/relationship"
	"github.com/sourcenetwork/acp_core/internal/theorem"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func HandleSimulateRequest(ctx context.Context, req *types.SimulateRequest) (*types.SimulateResponse, error) {
	manager, err := runtime.NewRuntimeManager(runtime.WithMemKV())
	if err != nil {
		return nil, newSimulateErr(err)
	}

	engine, err := zanzi.NewZanzi(manager.GetKVStore(), manager.GetLogger())
	if err != nil {
		return nil, newSimulateErr(err)
	}

	polHandler := policy.CreatePolicyHandler{}
	polResp, err := polHandler.Execute(ctx, manager, &types.CreatePolicyRequest{
		Policy:       req.Policy,
		MarshalType:  req.MarshalType,
		CreationTime: prototypes.TimestampNow(),
	})
	if err != nil {
		return nil, newSimulateErr(err)
	}
	polId := polResp.Policy.Id

	relationships, err := parser.ParseRelationships(req.RelationshipSet)
	if err != nil {
		return nil, newSimulateErr(err)
	}

	policyTheorem, err := parser.ParsePolicyTheorem(req.PolicyTheorem)
	if err != nil {
		return nil, newSimulateErr(err)
	}

	for _, rel := range relationships {
		var err error
		if rel.Relation == policy.OwnerRelation {
			principal, err := auth.NewDIDPrincipal(rel.Subject.GetActor().Id)
			if err != nil {
				return nil, err
			}
			ctx = auth.InjectPrincipal(ctx, principal)
			handler := relationship.RegisterObjectHandler{}
			_, err = handler.Execute(ctx, manager, &types.RegisterObjectRequest{
				PolicyId:     polId,
				Object:       rel.Object,
				CreationTime: prototypes.TimestampNow(),
			})
		} else {
			handler := relationship.SetRelationshipHandler{}
			_, err = handler.Execute(ctx, manager, &types.SetRelationshipRequest{
				PolicyId:     polId,
				CreationTime: prototypes.TimestampNow(),
				Relationship: rel,
			})
		}
		if err != nil {
			return nil, newSimulateErr(err)
		}
	}

	evaluator := theorem.NewEvaluator(engine)
	result, err := evaluator.EvaluatePolicyTheorem(ctx, polResp.Policy.Id, policyTheorem)
	if err != nil {
		return nil, newSimulateErr(err)
	}

	return &types.SimulateResponse{
		Policy:        polResp.Policy,
		Relationships: relationships,
		PolicyTheorem: policyTheorem,
		Result:        result,
	}, nil
}
