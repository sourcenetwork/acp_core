package sandbox

import (
	"context"
	"fmt"

	prototypes "github.com/cosmos/gogoproto/types"

	"github.com/sourcenetwork/acp_core/internal/parser"
	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/internal/raccoon"
	"github.com/sourcenetwork/acp_core/internal/relationship"
	"github.com/sourcenetwork/acp_core/internal/theorem"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/playground"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

func HandleNewSandboxRequest(ctx context.Context, manager runtime.RuntimeManager, req *playground.NewSandboxRequest) (*playground.NewSandboxResponse, error) {
	counter := raccoon.NewCounterStore(manager, sandboxStorePrefix)

	handle, err := counter.GetNext(ctx)
	if err != nil {
		return nil, err
	}

	name := req.Name
	if name == "" {
		name = fmt.Sprintf("%v", handle)
	}

	record := &playground.SandboxRecord{
		Name:        name,
		Handle:      handle,
		Description: req.Description,
	}

	repository := NewSandboxRepository(manager)
	err = repository.SetRecord(ctx, record)
	if err != nil {
		return nil, err
	}

	err = counter.Increment(ctx)
	if err != nil {
		return nil, err
	}

	return &playground.NewSandboxResponse{
		Record: record,
	}, nil
}

func HandleListSandboxes(ctx context.Context, manager runtime.RuntimeManager, req *playground.ListSandboxesRequest) (*playground.ListSandboxesResponse, error) {
	repository := NewSandboxRepository(manager)
	sandboxes, err := repository.ListSandboxes(ctx)
	if err != nil {
		return nil, err
	}
	return &playground.ListSandboxesResponse{
		Records: sandboxes,
	}, nil
}

type parsedCtx struct {
	PolicyDefinition string
	Relationships    []parser.LocatedObject[*types.Relationship]
	Theorem          *parser.LocatedPolicyTheorem
	Policy           *types.Policy
}

func (c *parsedCtx) ToCtx() *playground.SandboxCtx {
	return &playground.SandboxCtx{
		Policy:        c.Policy,
		Relationships: utils.MapSlice(c.Relationships, func(o parser.LocatedObject[*types.Relationship]) *types.Relationship { return o.Obj }),
		PolicyTheorem: c.Theorem.ToPolicyTheorem(),
	}
}

type SetStateHandler struct{}

func (h *SetStateHandler) Handle(ctx context.Context, manager runtime.RuntimeManager, req *playground.SetStateRequest) (*playground.SetStateResponse, error) {
	repository := NewSandboxRepository(manager)

	record, err := repository.GetSandbox(ctx, req.Handle)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, errors.Wrap("sandbox", errors.ErrorType_NOT_FOUND, errors.Pair("id", req.Handle))
	}

	record.Scratchpad = req.Data
	err = repository.SetRecord(ctx, record)
	if err != nil {
		return nil, err // TODO
	}

	simCtx, errs, err := h.parseCtx(ctx, manager, req.Data)
	if err != nil {
		return nil, err
	}
	if errs.HasErrors() {
		return &playground.SetStateResponse{
			Ok:     false,
			Errors: errs,
		}, nil
	}

	errs, err = h.populateEngine(ctx, manager, req.Handle, simCtx)
	if err != nil {
		return nil, err
	}
	if errs.HasErrors() {
		return &playground.SetStateResponse{
			Ok:     false,
			Errors: errs,
		}, nil
	}

	record.Data = req.Data
	record.Initialized = true
	record.Ctx = simCtx.ToCtx()
	err = repository.SetRecord(ctx, record)
	if err != nil {
		return nil, err
	}

	return &playground.SetStateResponse{
		Ok:     true,
		Errors: &playground.SandboxDataErrors{},
		Record: record,
	}, nil
}

func (h *SetStateHandler) parseCtx(ctx context.Context, manager runtime.RuntimeManager, data *playground.SandboxData) (*parsedCtx, *playground.SandboxDataErrors, error) {
	var errs = &playground.SandboxDataErrors{}

	// FIXME do full parsing once independent parsing is implemented
	_, err := policy.Unmarshal(data.PolicyDefinition, types.PolicyMarshalingType_SHORT_YAML)
	if err != nil {
		err := &errors.ParserMessage{
			Message:   err.Error(),
			Sevirity:  errors.Severity_ERROR,
			InputName: "policy",
			// Range is empty because unmarshaling still doesn't support that feature
		}
		errs.PolicyErrors = append(errs.PolicyErrors, err)
	}

	relationships, err := parser.ParseRelationshipsWithLocation(data.Relationships)
	if err != nil {
		if parserErr, ok := err.(*errors.ParserReport); ok {
			errs.RelationshipsErrors = append(errs.RelationshipsErrors, parserErr.Messages...)
		} else {
			return nil, nil, err // TODO WRAP
		}
	}

	theorem, err := parser.ParsePolicyTheorem(data.PolicyTheorem)
	if err != nil {
		if parserErr, ok := err.(*errors.ParserReport); ok {
			errs.TheoremsErrrors = append(errs.TheoremsErrrors, parserErr.Messages...)
		} else {
			return nil, nil, err // TODO
		}
	}
	if errs.HasErrors() {
		return nil, errs, nil
	}

	simCtx := &parsedCtx{
		Relationships:    relationships,
		PolicyDefinition: data.PolicyDefinition,
		Theorem:          theorem,
	}
	return simCtx, &playground.SandboxDataErrors{}, nil
}

func (h *SetStateHandler) populateEngine(ctx context.Context, manager runtime.RuntimeManager, handle uint64, simCtx *parsedCtx) (*playground.SandboxDataErrors, error) {
	var errs = &playground.SandboxDataErrors{}

	engineManager, err := GetManagerForSandbox(manager, handle)
	if err != nil {
		return nil, err
	}

	polHandler := policy.CreatePolicyHandler{}
	polResp, err := polHandler.Execute(ctx, engineManager, &types.CreatePolicyRequest{
		Policy:       simCtx.PolicyDefinition,
		MarshalType:  types.PolicyMarshalingType_SHORT_YAML,
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
			errs.PolicyErrors = append(errs.PolicyErrors, err)
			return errs, nil
		} else {
			// non marshaling errors should terminate execution
			return nil, err // TODO figure out error type for the result
		}
	}
	simCtx.Policy = polResp.Policy

	return h.setRelationships(ctx, engineManager, simCtx)
}

/*
func (h *SetStateHandler) validateRelationships(ctx context.Context, manager runtime.RuntimeManager, simCtx *parsedCtx) (*playground.SandboxDataErrors, error) {
	var declarationErrors = &playground.SandboxDataErrors{}

	for _, indexedRel := range simCtx.Relationships {
		valid, errMsg, err := relationship.ValidateRelationship(ctx, manager, simCtx.Policy.Id, indexedRel.Obj)
		if err != nil {
			return nil, err // TODO
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
*/

func (h *SetStateHandler) setRelationships(ctx context.Context, manager runtime.RuntimeManager, simCtx *parsedCtx) (*playground.SandboxDataErrors, error) {
	errs := &playground.SandboxDataErrors{}
	ownerLookup := make(map[string]auth.Principal)
	ownerRels, rels := utils.PartitionSlice(simCtx.Relationships, func(obj parser.LocatedObject[*types.Relationship]) bool {
		return obj.Obj.Relation == policy.OwnerRelation
	})

	for _, obj := range ownerRels {
		principal, err := auth.NewDIDPrincipal(obj.Obj.Subject.GetActor().Id)
		if err != nil {
			return nil, err
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
		ownerLookup[obj.Obj.Object.String()] = principal
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

func HandleVerifyTheorem(ctx context.Context, manager runtime.RuntimeManager, req *playground.VerifyTheoremsRequest) (*playground.VerifyTheoremsResponse, error) {
	repository := NewSandboxRepository(manager)

	record, err := repository.GetSandbox(ctx, req.Handle)
	if err != nil {
		return nil, err // TODO WRAP
	}
	if record == nil {
		return nil, errors.Wrap("sandbox not found", errors.ErrorType_NOT_FOUND, errors.Pair("handle", req.Handle))
	}
	if !record.Initialized {
		return nil, errors.Wrap("uninitialized sandbox cannot execute theorems", errors.ErrorType_OPERATION_FORBIDDEN, errors.Pair("handle", req.Handle))
	}

	manager, err = GetManagerForSandbox(manager, req.Handle)
	if err != nil {
		return nil, err // TODO WRAP
	}

	engine, err := zanzi.NewZanzi(manager.GetKVStore(), manager.GetLogger())
	if err != nil {
		return nil, err // TODO
	}
	evaluator := theorem.NewEvaluator(engine)

	result, err := evaluator.EvaluatePolicyTheoremDSL(ctx, record.Ctx.Policy.Id, record.Data.PolicyTheorem)
	if err != nil {
		return nil, err
	}

	return &playground.VerifyTheoremsResponse{
		Result: result,
	}, nil
}

func HandleGetCatalogue(ctx context.Context, manager runtime.RuntimeManager, req *playground.GetCatalogueRequest) (*playground.GetCatalogueResponse, error) {
	repository := NewSandboxRepository(manager)

	record, err := repository.GetSandbox(ctx, req.Handle)
	if err != nil {
		return nil, err // TODO WRAP
	}
	if record == nil {
		return nil, errors.Wrap("sandbox not found", errors.ErrorType_NOT_FOUND, errors.Pair("handle", req.Handle))
	}
	if !record.Initialized {
		return nil, errors.Wrap("uninitialized sandbox cannot execute theorems", errors.ErrorType_OPERATION_FORBIDDEN, errors.Pair("handle", req.Handle))
	}

	manager, err = GetManagerForSandbox(manager, req.Handle)
	if err != nil {
		return nil, err // TODO WRAP
	}

	catalogue, err := policy.HandleBuildCatalogue(ctx, manager, record.Ctx.Policy.Id)
	if err != nil {
		return nil, err // TODO wrap
	}

	return &playground.GetCatalogueResponse{
		Catalogue: catalogue,
	}, nil
}

func HandleRestoreScratchpad(ctx context.Context, manager runtime.RuntimeManager, req *playground.RestoreScratchpadRequest) (*playground.RestoreScratchpadResponse, error) {
	repository := NewSandboxRepository(manager)

	record, err := repository.GetSandbox(ctx, req.Handle)
	if err != nil {
		return nil, err // TODO WRAP
	}
	if record == nil {
		return nil, errors.Wrap("sandbox not found", errors.ErrorType_NOT_FOUND, errors.Pair("handle", req.Handle))
	}
	if !record.Initialized {
		return nil, errors.Wrap("uninitialized sandbox cannot execute theorems", errors.ErrorType_OPERATION_FORBIDDEN, errors.Pair("handle", req.Handle))
	}

	record.Scratchpad = record.Data

	err = repository.SetRecord(ctx, record)
	if err != nil {
		return nil, err
	}

	return &playground.RestoreScratchpadResponse{
		Scratchpad: record.Scratchpad,
	}, nil
}
