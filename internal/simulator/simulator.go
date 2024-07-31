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
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

func SimulateDeclaration(ctx context.Context, manager runtime.RuntimeManager, declaration *types.SimulationCtxDeclaration) (*types.AnnotatedSimulationResult, error) {
	parsedCtx, parseErrs, err := parseDeclaration(ctx, manager, declaration)
	if err != nil {
		return nil, err
	}
	if parseErrs.HasErrors() {
		return &types.AnnotatedSimulationResult{
			Declaration:         declaration,
			Ctx:                 nil,
			Errors:              parseErrs,
			PolicyTheoremResult: nil,
		}, nil
	}

	engine, err := zanzi.NewZanzi(manager.GetKVStore(), manager.GetLogger())
	if err != nil {
		return nil, newSimulateErr(err)
	}

	simCtx := parsedCtx.ToCtx()

	parseErrs, err = validateRelationships(ctx, manager, parsedCtx)
	if err != nil {
		return nil, newSimulateErr(err)
	}
	if parseErrs.HasErrors() {
		return &types.AnnotatedSimulationResult{
			Declaration:         declaration,
			Ctx:                 simCtx,
			Errors:              parseErrs,
			PolicyTheoremResult: nil,
		}, nil
	}

	parseErrs, err = setRelationships(ctx, manager, parsedCtx)
	if err != nil {
		return nil, newSimulateErr(err)
	}
	if parseErrs.HasErrors() {
		return &types.AnnotatedSimulationResult{
			Declaration:         declaration,
			Ctx:                 simCtx,
			Errors:              parseErrs,
			PolicyTheoremResult: nil,
		}, nil
	}

	evaluator := theorem.NewEvaluator(engine)
	annotatedResult, err := evaluator.EvaluatePolicyTheoremDSL(ctx, simCtx.Policy.Id, declaration.PolicyTheorem)
	if err != nil {
		// since the theorem was previously parsed
		// any resulting errors won't be parse errors, therefore we can early terminate here
		return nil, newSimulateErr(err)
	}

	return &types.AnnotatedSimulationResult{
		Declaration:         declaration,
		Ctx:                 simCtx,
		Errors:              nil,
		PolicyTheoremResult: annotatedResult,
	}, nil
}

// parsedCtx is a container type containing the result of parsing a SimulationCtxDeclaration
type parsedCtx struct {
	Policy        *types.Policy
	Relationships []parser.IndexedObject[*types.Relationship]
	Theorem       *parser.IndexedPolicyTheorem
}

func (c *parsedCtx) ToCtx() *types.SimulationCtx {
	return &types.SimulationCtx{
		Policy:        c.Policy,
		Relationships: utils.MapSlice(c.Relationships, func(o parser.IndexedObject[*types.Relationship]) *types.Relationship { return o.Obj }),
		PolicyTheorem: c.Theorem.ToPolicyTheorem(),
	}
}

// parseDeclaration processes a Ctx Declaration and returns the concrete parsed types.
// Non fatal input errors found from parsing the declaration are added to DeclarationErrors.
// Any fatal error encountered in the process is returned as an error
func parseDeclaration(ctx context.Context, manager runtime.RuntimeManager, declaration *types.SimulationCtxDeclaration) (*parsedCtx, *types.DeclarationErrors, error) {
	var declarationErrors = &types.DeclarationErrors{}

	// FIXME ideally I would be able to parse the Policy independently, but for now I need to parse and create in one step
	polHandler := policy.CreatePolicyHandler{}
	polResp, err := polHandler.Execute(ctx, manager, &types.CreatePolicyRequest{
		Policy:       declaration.Policy,
		MarshalType:  declaration.MarshalType,
		CreationTime: prototypes.TimestampNow(),
	})
	if err != nil {
		if errors.Is(err, policy.ErrUnmarshaling) {
			err := &errors.ParserMessage{
				Message:   err.Error(),
				Sevirity:  errors.Severity_ERROR,
				InputName: "policy",
				// Range is empty because unmarshaling still doesn't support that feature
			}
			declarationErrors.PolicyErrors = append(declarationErrors.PolicyErrors, err)
		} else {
			// non marshaling errors should terminate execution
			return nil, nil, newSimulateErr(err) // TODO figure out error type for the result
		}
	}

	relationships, err := parser.ParseRelationshipsWithPosition(declaration.RelationshipSet)
	if err != nil {
		if parserErr, ok := err.(*errors.ParserReport); ok {
			declarationErrors.RelationshipsErrors = append(declarationErrors.RelationshipsErrors, parserErr.Messages...)
		} else {
			return nil, nil, newSimulateErr(err)
		}
	}

	theorem, err := parser.ParsePolicyTheorem(declaration.PolicyTheorem)
	if err != nil {
		if parserErr, ok := err.(*errors.ParserReport); ok {
			declarationErrors.TheoremsErrrors = append(declarationErrors.TheoremsErrrors, parserErr.Messages...)
		} else {
			return nil, nil, newSimulateErr(err)
		}
	}

	if declarationErrors.HasErrors() {
		return nil, declarationErrors, nil
	}

	simCtx := &parsedCtx{
		Relationships: relationships,
		Policy:        polResp.Policy,
		Theorem:       theorem,
	}
	return simCtx, declarationErrors, nil
}

func validateRelationships(ctx context.Context, manager runtime.RuntimeManager, simCtx *parsedCtx) (*types.DeclarationErrors, error) {
	var declarationErrors = &types.DeclarationErrors{}

	for _, indexedRel := range simCtx.Relationships {
		valid, errMsg, err := relationship.ValidateRelationship(ctx, manager, simCtx.Policy.Id, indexedRel.Obj)
		if err != nil {
			return nil, newSimulateErr(err)
		}
		if !valid {
			msg := &errors.ParserMessage{
				Message:   errMsg,
				Sevirity:  errors.Severity_ERROR,
				InputName: indexedRel.Obj.String(),
				Range: &errors.BufferRange{
					Start: &errors.BufferPosition{
						Column: indexedRel.Range.Start.Column,
						Line:   indexedRel.Range.Start.Line,
					},
					End: &errors.BufferPosition{
						Column: indexedRel.Range.End.Column,
						Line:   indexedRel.Range.End.Line,
					},
				},
			}
			declarationErrors.RelationshipsErrors = append(declarationErrors.RelationshipsErrors, msg)
		}
	}
	return declarationErrors, nil
}

func setRelationships(ctx context.Context, manager runtime.RuntimeManager, simCtx *parsedCtx) (*types.DeclarationErrors, error) {
	errs := &types.DeclarationErrors{}
	ownerLookup := make(map[string]auth.Principal)
	ownerRels, rels := utils.PartitionSlice(simCtx.Relationships, func(obj parser.IndexedObject[*types.Relationship]) bool {
		return obj.Obj.Relation == policy.OwnerRelation
	})

	for _, obj := range ownerRels {
		principal, err := auth.NewDIDPrincipal(obj.Obj.Subject.GetActor().Id)
		if err != nil {
			return nil, newSimulateErr(err)
		}
		authenticatedCtx := auth.InjectPrincipal(ctx, principal)
		handler := relationship.RegisterObjectHandler{}
		_, err = handler.Execute(authenticatedCtx, manager, &types.RegisterObjectRequest{
			PolicyId:     simCtx.Policy.Id,
			Object:       obj.Obj.Object,
			CreationTime: prototypes.TimestampNow(),
		})
		if err != nil {
			return nil, err // TODO
		}
		ownerLookup[obj.Obj.String()] = principal
	}

	for _, indexedObj := range rels {
		principal, ok := ownerLookup[indexedObj.Obj.Object.String()]
		if !ok {
			msg := &errors.ParserMessage{
				Message:   "object not registered",
				Sevirity:  errors.Severity_ERROR,
				InputName: "",
				Range: &errors.BufferRange{
					Start: &errors.BufferPosition{
						Column: indexedObj.Range.Start.Column,
						Line:   indexedObj.Range.Start.Line,
					},
					End: &errors.BufferPosition{
						Column: indexedObj.Range.End.Column,
						Line:   indexedObj.Range.End.Line,
					},
				},
			}
			errs.RelationshipsErrors = append(errs.RelationshipsErrors, msg)
			continue
		}
		authenticatedCtx := auth.InjectPrincipal(ctx, principal)
		handler := relationship.SetRelationshipHandler{}
		_, err := handler.Execute(authenticatedCtx, manager, &types.SetRelationshipRequest{
			PolicyId:     simCtx.Policy.Id,
			CreationTime: prototypes.TimestampNow(),
			Relationship: indexedObj.Obj,
		})
		if err != nil {
			return nil, err // TODO figure out err, rels should have been validated so this would be ok?
		}

	}
	return errs, nil
}
