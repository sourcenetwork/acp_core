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
	counter := raccoon.NewCounterStoreFromRunetimeManager(manager, sandboxStorePrefix)

	releaser := counter.Acquire()
	defer releaser.Release()

	handle, err := counter.GetNext(ctx)
	if err != nil {
		return nil, newNewSandboxErr(err)
	}

	if req.Name == "" {
		req.Name = fmt.Sprintf("%v", handle)
	}

	record := &playground.SandboxRecord{
		Name:        req.Name,
		Handle:      handle,
		Description: req.Description,
	}

	repository := NewSandboxRepository(manager.GetKVStore())
	err = repository.SetRecord(ctx, record)
	if err != nil {
		return nil, newNewSandboxErr(err)
	}

	err = counter.Increment(ctx)
	if err != nil {
		return nil, newNewSandboxErr(err)
	}

	return &playground.NewSandboxResponse{
		Record: record,
	}, nil
}

func HandleListSandboxes(ctx context.Context, manager runtime.RuntimeManager, req *playground.ListSandboxesRequest) (*playground.ListSandboxesResponse, error) {
	repository := NewSandboxRepository(manager.GetKVStore())
	sandboxes, err := repository.ListSandboxes(ctx)
	if err != nil {
		return nil, newListSandboxesErr(err)
	}

	return &playground.ListSandboxesResponse{
		Records: sandboxes,
	}, nil
}

type SetStateHandler struct{}

func (h *SetStateHandler) Handle(ctx context.Context, manager runtime.RuntimeManager, req *playground.SetStateRequest) (*playground.SetStateResponse, error) {
	repository := NewSandboxRepository(manager.GetKVStore())

	record, err := repository.GetSandbox(ctx, req.Handle)
	if err != nil {
		return nil, newSetStateErr(err, req.Handle)
	}
	if record == nil {
		return nil, newSetStateErr(errors.Wrap("sandbox", errors.ErrorType_NOT_FOUND), req.Handle)
	}

	record.Scratchpad = req.Data
	err = repository.SetRecord(ctx, record)
	if err != nil {
		return nil, newSetStateErr(err, req.Handle)
	}

	simCtx, errs, err := h.parseCtx(ctx, manager, req.Data)
	if err != nil {
		return nil, newSetStateErr(err, req.Handle)
	}
	if errs.HasErrors() {
		return &playground.SetStateResponse{
			Ok:     false,
			Errors: errs,
		}, nil
	}

	sandboxManager, err := GetManagerForSandbox(manager, req.Handle)
	if err != nil {
		return nil, err
	}

	errs, err = h.populateEngine(ctx, sandboxManager, req.Handle, simCtx)
	if err != nil {
		return nil, newSetStateErr(err, req.Handle)
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
		return nil, newSetStateErr(err, req.Handle)
	}

	return &playground.SetStateResponse{
		Ok:     true,
		Errors: &playground.SandboxDataErrors{},
		Record: record,
	}, nil
}

// parseCtx parses the input data and returns a parsed ctx or all errors found while parsing or
// any other errors encountered during the program execution
func (h *SetStateHandler) parseCtx(ctx context.Context, manager runtime.RuntimeManager, data *playground.SandboxData) (*parsedSandboxCtx, *playground.SandboxDataErrors, error) {
	var errs = &playground.SandboxDataErrors{}

	// FIXME do full parsing once independent parsing is implemented
	_, err := policy.Unmarshal(data.PolicyDefinition, types.PolicyMarshalingType_SHORT_YAML)
	if err != nil {
		err := &types.LocatedMessage{
			Message:   err.Error(),
			Kind:      types.LocatedMessage_ERROR,
			InputName: "policy",
			// Range is empty because unmarshaling still doesn't support that feature
		}
		errs.PolicyErrors = append(errs.PolicyErrors, err)
	}

	relationships, report := parser.ParseRelationshipsWithLocation(data.Relationships)
	errs.RelationshipsErrors = append(errs.RelationshipsErrors, report.GetMessages()...)

	theorem, report := parser.ParsePolicyTheorem(data.PolicyTheorem)
	errs.TheoremsErrors = append(errs.TheoremsErrors, report.GetMessages()...)

	if errs.HasErrors() {
		return nil, errs, nil
	}

	simCtx := &parsedSandboxCtx{
		Relationships:    relationships,
		PolicyDefinition: data.PolicyDefinition,
		Theorem:          theorem,
	}
	return simCtx, &playground.SandboxDataErrors{}, nil
}

func (h *SetStateHandler) populateEngine(ctx context.Context, manager runtime.RuntimeManager, handle uint64, simCtx *parsedSandboxCtx) (*playground.SandboxDataErrors, error) {
	errs, err := h.setPolicy(ctx, manager, handle, simCtx)
	if err != nil {
		return nil, err
	}
	if errs.HasErrors() {
		return errs, nil
	}

	ownerMap, errs, err := h.registerObjects(ctx, manager, handle, simCtx)
	if err != nil {
		return nil, err
	}

	errs2, err := h.setRelationships(ctx, manager, simCtx, ownerMap)
	if err != nil {
		return nil, err
	}

	errs2.PolicyErrors = append(errs2.PolicyErrors, errs.PolicyErrors...)
	errs2.RelationshipsErrors = append(errs2.RelationshipsErrors, errs.RelationshipsErrors...)
	errs2.TheoremsErrors = append(errs2.TheoremsErrors, errs.TheoremsErrors...)

	return errs2, nil
}

func (h *SetStateHandler) setPolicy(ctx context.Context, manager runtime.RuntimeManager, handle uint64, simCtx *parsedSandboxCtx) (*playground.SandboxDataErrors, error) {
	errs := &playground.SandboxDataErrors{}

	polHandler := policy.CreatePolicyHandler{}
	polResp, err := polHandler.Execute(ctx, manager, &types.CreatePolicyRequest{
		Policy:       simCtx.PolicyDefinition,
		MarshalType:  types.PolicyMarshalingType_SHORT_YAML,
		CreationTime: prototypes.TimestampNow(),
	})

	if err != nil {
		if errors.Is(err, errors.ErrorType_BAD_INPUT) {
			msg := &types.LocatedMessage{
				Message:   err.Error(),
				Kind:      types.LocatedMessage_ERROR,
				InputName: "policy",
				// Range is empty because unmarshaling still doesn't support that feature
			}
			errs.PolicyErrors = append(errs.PolicyErrors, msg)
		} else {
			// non marshaling errors should terminate execution
			return nil, err
		}
	}

	simCtx.Policy = polResp.Policy
	return errs, nil
}

func (h *SetStateHandler) registerObjects(ctx context.Context, manager runtime.RuntimeManager, handle uint64, simCtx *parsedSandboxCtx) (map[string]auth.Principal, *playground.SandboxDataErrors, error) {
	errs := &playground.SandboxDataErrors{}
	ownerLookup := make(map[string]auth.Principal)

	ownerRels := utils.FilterSlice(simCtx.Relationships, func(obj parser.LocatedObject[*types.Relationship]) bool {
		return obj.Obj.Relation == policy.OwnerRelation
	})

	for _, obj := range ownerRels {
		principal, err := auth.NewDIDPrincipal(obj.Obj.Subject.GetActor().Id)
		// creating a principal should only fail if the actor is invalid
		// meaning the relationship is invalid
		if err != nil {
			err := &types.LocatedMessage{
				Message:   err.Error(),
				Kind:      types.LocatedMessage_ERROR,
				InputName: "relationships",
				Range:     obj.Range,
			}
			errs.RelationshipsErrors = append(errs.RelationshipsErrors, err)
			continue
		}

		authenticatedCtx := auth.InjectPrincipal(ctx, principal)

		handler := relationship.RegisterObjectHandler{}
		_, err = handler.Execute(authenticatedCtx, manager, &types.RegisterObjectRequest{
			PolicyId:     simCtx.Policy.Id,
			Object:       obj.Obj.Object,
			CreationTime: prototypes.TimestampNow(),
		})
		if err != nil {
			if errors.Is(err, errors.ErrorType_BAD_INPUT) {
				err := &types.LocatedMessage{
					Message:   err.Error(),
					Kind:      types.LocatedMessage_ERROR,
					InputName: "relationships",
					Range:     obj.Range,
				}
				errs.RelationshipsErrors = append(errs.RelationshipsErrors, err)
			} else {
				return nil, nil, err
			}
		}
		ownerLookup[obj.Obj.Object.String()] = principal
	}
	return ownerLookup, errs, nil
}

func (h *SetStateHandler) setRelationships(ctx context.Context, manager runtime.RuntimeManager, simCtx *parsedSandboxCtx, ownerMap map[string]auth.Principal) (*playground.SandboxDataErrors, error) {
	errs := &playground.SandboxDataErrors{}
	rels := utils.FilterSlice(simCtx.Relationships, func(obj parser.LocatedObject[*types.Relationship]) bool {
		return obj.Obj.Relation != policy.OwnerRelation
	})

	for _, indexedObj := range rels {
		principal, ok := ownerMap[indexedObj.Obj.Object.String()]
		if !ok {
			msg := &types.LocatedMessage{
				Message:   "object not registered",
				Kind:      types.LocatedMessage_ERROR,
				InputName: "relationships",
				Range:     indexedObj.Range,
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
			if errors.Is(err, errors.ErrorType_BAD_INPUT) {
				err := &types.LocatedMessage{
					Message:   err.Error(),
					Kind:      types.LocatedMessage_ERROR,
					InputName: "relationships",
					Range:     indexedObj.Range,
				}
				errs.RelationshipsErrors = append(errs.RelationshipsErrors, err)
			} else {
				return nil, err
			}
		}
	}
	return errs, nil
}

func HandleVerifyTheorem(ctx context.Context, manager runtime.RuntimeManager, req *playground.VerifyTheoremsRequest) (*playground.VerifyTheoremsResponse, error) {
	repository := NewSandboxRepository(manager.GetKVStore())

	record, err := repository.GetSandbox(ctx, req.Handle)
	if err != nil {
		return nil, newVerifyTheoremsErr(err, req.Handle)
	}
	if record == nil {
		return nil, errors.Wrap("sandbox not found", errors.ErrorType_NOT_FOUND, errors.Pair("handle", req.Handle))
	}
	if !record.Initialized {
		return nil, errors.Wrap("uninitialized sandbox cannot execute theorems", errors.ErrorType_OPERATION_FORBIDDEN, errors.Pair("handle", req.Handle))
	}

	manager, err = GetManagerForSandbox(manager, req.Handle)
	if err != nil {
		return nil, newVerifyTheoremsErr(err, req.Handle)
	}

	engine, err := zanzi.NewZanzi(manager.GetKVStore(), manager.GetLogger())
	if err != nil {
		return nil, newVerifyTheoremsErr(err, req.Handle)
	}
	evaluator := theorem.NewEvaluator(engine)

	result, err := evaluator.EvaluatePolicyTheoremDSL(ctx, record.Ctx.Policy.Id, record.Data.PolicyTheorem)
	if err != nil {
		return nil, newVerifyTheoremsErr(err, req.Handle)
	}

	return &playground.VerifyTheoremsResponse{
		Result: result,
	}, nil
}

func HandleGetCatalogue(ctx context.Context, manager runtime.RuntimeManager, req *playground.GetCatalogueRequest) (*playground.GetCatalogueResponse, error) {
	repository := NewSandboxRepository(manager.GetKVStore())

	record, err := repository.GetSandbox(ctx, req.Handle)
	if err != nil {
		return nil, newGetCatalogueErr(err, req.Handle)
	}
	if record == nil {
		return nil, errors.Wrap("sandbox not found", errors.ErrorType_NOT_FOUND, errors.Pair("handle", req.Handle))
	}
	if !record.Initialized {
		err := errors.Wrap("uninitialized sandbox cannot execute theorems",
			errors.ErrorType_OPERATION_FORBIDDEN)
		return nil, newGetCatalogueErr(err, req.Handle)
	}

	manager, err = GetManagerForSandbox(manager, req.Handle)
	if err != nil {
		return nil, newGetCatalogueErr(err, req.Handle)
	}

	engine, err := zanzi.NewZanzi(manager.GetKVStore(), manager.GetLogger())
	if err != nil {
		return nil, newGetCatalogueErr(err, req.Handle)
	}
	catalogue, err := policy.BuildCatalogue(ctx, engine, record.Ctx.Policy.Id)
	if err != nil {
		return nil, newGetCatalogueErr(err, req.Handle)
	}

	return &playground.GetCatalogueResponse{
		Catalogue: catalogue,
	}, nil
}

func HandleRestoreScratchpad(ctx context.Context, manager runtime.RuntimeManager, req *playground.RestoreScratchpadRequest) (*playground.RestoreScratchpadResponse, error) {
	repository := NewSandboxRepository(manager.GetKVStore())

	record, err := repository.GetSandbox(ctx, req.Handle)
	if err != nil {
		return nil, newRestoreScratchpadErr(err, req.Handle)
	}
	if record == nil {
		err = errors.Wrap("sandbox not found", errors.ErrorType_NOT_FOUND)
		return nil, newRestoreScratchpadErr(err, req.Handle)
	}
	if !record.Initialized {
		err = errors.Wrap("uninitialized sandbox cannot execute theorems", errors.ErrorType_OPERATION_FORBIDDEN)
		return nil, newRestoreScratchpadErr(err, req.Handle)
	}

	record.Scratchpad = record.Data

	err = repository.SetRecord(ctx, record)
	if err != nil {
		return nil, newRestoreScratchpadErr(err, req.Handle)
	}

	return &playground.RestoreScratchpadResponse{
		Scratchpad: record.Scratchpad,
	}, nil
}
