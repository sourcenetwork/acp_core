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

func SimulateDeclaration(ctx context.Context, manager runtime.RuntimeManager, declaration *types.SimulationCtxDeclaration) (*types.AnnotatedSimulationResult, error) {
	engine, err := zanzi.NewZanzi(manager.GetKVStore(), manager.GetLogger())
	if err != nil {
		return nil, newSimulateErr(err)
	}

	polHandler := policy.CreatePolicyHandler{}
	polResp, err := polHandler.Execute(ctx, manager, &types.CreatePolicyRequest{
		Policy:       declaration.Policy,
		MarshalType:  declaration.MarshalType,
		CreationTime: prototypes.TimestampNow(),
	})
	if err != nil {
		return nil, newSimulateErr(err) // TODO figure out error type for the result
	}
	polId := polResp.Policy.Id

	relationships, err := parser.ParseRelationships(declaration.RelationshipSet)
	if err != nil {
		return nil, newSimulateErr(err) // TODO figure out error types
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
	annotatedResult, err := evaluator.EvaluatePolicyTheoremDSL(ctx, polResp.Policy.Id, declaration.PolicyTheorem)
	if err != nil {
		return nil, newSimulateErr(err) // TODO figure out error type
	}

	return &types.AnnotatedSimulationResult{
		Ctx: &types.SimulationCtx{
			Policy:        polResp.Policy,
			Relationships: relationships,
			PolicyTheorem: annotatedResult.Theorem,
		},
		Errors:              nil,
		PolicyTheoremResult: annotatedResult,
	}, nil
}
