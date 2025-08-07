package sandbox

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/internal/raccoon"
	"github.com/sourcenetwork/acp_core/internal/relationship"
	"github.com/sourcenetwork/acp_core/internal/theorem"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/parser"
	"github.com/sourcenetwork/acp_core/pkg/parser/theorem_parser"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

func HandleNewSandboxRequest(ctx context.Context, manager runtime.RuntimeManager, req *types.NewSandboxRequest) (*types.NewSandboxResponse, error) {
	counter := raccoon.NewCounterStoreFromRuntimeManager(manager, sandboxCounterPrefix)

	releaser := counter.Acquire()
	defer releaser.Release()

	handle, err := counter.GetNext(ctx)
	if err != nil {
		return nil, newNewSandboxErr(err)
	}

	if req.Name == "" {
		req.Name = fmt.Sprintf("%v", handle)
	}

	record := &types.SandboxRecord{
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

	return &types.NewSandboxResponse{
		Record: record,
	}, nil
}

func HandleListSandboxes(ctx context.Context, manager runtime.RuntimeManager, req *types.ListSandboxesRequest) (*types.ListSandboxesResponse, error) {
	// FIXME: there's something weird going on in raccoon and I have no idea what.
	// Listing isn't working, but indiidual fetching does.
	// Workaround for now, I'm sorry.
	counter := raccoon.NewCounterStoreFromRuntimeManager(manager, sandboxCounterPrefix)
	max, err := counter.GetNext(ctx)
	if err != nil {
		return nil, newListSandboxesErr(err)
	}

	var records []*types.SandboxRecord
	repository := NewSandboxRepository(manager.GetKVStore())
	for i := uint64(1); i < max; i += 1 {
		record, err := repository.GetSandbox(ctx, i)
		if err != nil {
			return nil, newListSandboxesErr(err)
		}
		if record != nil {
			records = append(records, record)
		}
	}
	return &types.ListSandboxesResponse{
		Records: records,
	}, nil
}

type SetStateHandler struct{}

func (h *SetStateHandler) Handle(ctx context.Context, manager runtime.RuntimeManager, req *types.SetStateRequest) (*types.SetStateResponse, error) {
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
		return &types.SetStateResponse{
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
		return &types.SetStateResponse{
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

	return &types.SetStateResponse{
		Ok:     true,
		Errors: &types.SandboxDataErrors{},
		Record: record,
	}, nil
}

// parseCtx parses the input data and returns a parsed ctx or all errors found while parsing or
// any other errors encountered during the program execution
func (h *SetStateHandler) parseCtx(ctx context.Context, manager runtime.RuntimeManager, data *types.SandboxData) (*parsedSandboxCtx, *types.SandboxDataErrors, error) {
	var errs = &types.SandboxDataErrors{}

	// FIXME do full parsing once independent parsing is implemented
	_, err := policy.Unmarshal(data.PolicyDefinition, types.PolicyMarshalingType_SHORT_YAML)
	if err != nil {
		end := getPolicyEndPosition(data.PolicyDefinition)
		err := &types.LocatedMessage{
			Message:   err.Error(),
			Kind:      types.LocatedMessage_ERROR,
			InputName: "policy",
			Interval: &types.BufferInterval{
				Start: &types.BufferPosition{
					Line:   1,
					Column: 1,
				},
				End: &end,
			},
		}
		errs.PolicyErrors = append(errs.PolicyErrors, err)
	}

	relationships, report := theorem_parser.ParseRelationshipsWithLocation(data.Relationships)
	errs.RelationshipsErrors = append(errs.RelationshipsErrors, report.GetMessages()...)

	theorem, report := theorem_parser.ParsePolicyTheorem(data.PolicyTheorem)
	errs.TheoremsErrors = append(errs.TheoremsErrors, report.GetMessages()...)

	if errs.HasErrors() {
		return nil, errs, nil
	}

	simCtx := &parsedSandboxCtx{
		Relationships:    relationships,
		PolicyDefinition: data.PolicyDefinition,
		Theorem:          theorem,
	}
	return simCtx, &types.SandboxDataErrors{}, nil
}

func (h *SetStateHandler) populateEngine(ctx context.Context, manager runtime.RuntimeManager, handle uint64, simCtx *parsedSandboxCtx) (*types.SandboxDataErrors, error) {
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

func (h *SetStateHandler) setPolicy(ctx context.Context, manager runtime.RuntimeManager, handle uint64, simCtx *parsedSandboxCtx) (*types.SandboxDataErrors, error) {
	errs := &types.SandboxDataErrors{}

	authenticatedCtx := auth.InjectPrincipal(ctx, types.RootPrincipal())
	polHandler := policy.CreatePolicyHandler{}
	polResp, err := polHandler.Execute(authenticatedCtx, manager, &types.CreatePolicyRequest{
		Policy:      simCtx.PolicyDefinition,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	})

	if err != nil {
		end := getPolicyEndPosition(simCtx.PolicyDefinition)
		if errors.Is(err, errors.ErrorType_BAD_INPUT) {
			msg := &types.LocatedMessage{
				Message:   err.Error(),
				Kind:      types.LocatedMessage_ERROR,
				InputName: "policy",
				Interval: &types.BufferInterval{
					Start: &types.BufferPosition{
						Line:   1,
						Column: 1,
					},
					End: &end,
				},
			}
			errs.PolicyErrors = append(errs.PolicyErrors, msg)
			return errs, nil
		} else {
			// non marshaling errors should terminate execution
			return nil, err
		}
	}

	simCtx.Policy = polResp.Record.Policy
	return errs, nil
}

func (h *SetStateHandler) registerObjects(ctx context.Context, manager runtime.RuntimeManager, handle uint64, simCtx *parsedSandboxCtx) (map[string]types.Principal, *types.SandboxDataErrors, error) {
	errs := &types.SandboxDataErrors{}
	ownerLookup := make(map[string]types.Principal)

	ownerRels := utils.FilterSlice(simCtx.Relationships, func(obj parser.LocatedObject[*types.Relationship]) bool {
		return obj.Obj.Relation == policy.OwnerRelation
	})

	for _, obj := range ownerRels {
		// direct owner relationships are only defined for
		// actor subjects
		if obj.Obj.Subject.GetActor() == nil {
			err := &types.LocatedMessage{
				Message:   "invalid relationship: invalid subject: owner relationship requires a `did` actor, make sure actor is a did",
				Kind:      types.LocatedMessage_ERROR,
				InputName: "relationships",
				Interval:  obj.Interval,
			}
			errs.RelationshipsErrors = append(errs.RelationshipsErrors, err)
			continue
		}

		// creating a principal should only fail if the actor is invalid
		// meaning the relationship is invalid
		principal, err := types.NewDIDPrincipal(obj.Obj.Subject.GetActor().Id)
		if err != nil {
			err := &types.LocatedMessage{
				Message:   err.Error(),
				Kind:      types.LocatedMessage_ERROR,
				InputName: "relationships",
				Interval:  obj.Interval,
			}
			errs.RelationshipsErrors = append(errs.RelationshipsErrors, err)
			continue
		}

		authenticatedCtx := auth.InjectPrincipal(ctx, principal)

		handler := relationship.RegisterObjectHandler{}
		_, err = handler.Execute(authenticatedCtx, manager, &types.RegisterObjectRequest{
			PolicyId: simCtx.Policy.Id,
			Object:   obj.Obj.Object,
		})
		if err != nil {
			if errors.Is(err, errors.ErrorType_BAD_INPUT) {
				err := &types.LocatedMessage{
					Message:   err.Error(),
					Kind:      types.LocatedMessage_ERROR,
					InputName: "relationships",
					Interval:  obj.Interval,
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

func (h *SetStateHandler) setRelationships(ctx context.Context, manager runtime.RuntimeManager, simCtx *parsedSandboxCtx, ownerMap map[string]types.Principal) (*types.SandboxDataErrors, error) {
	errs := &types.SandboxDataErrors{}
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
				Interval:  indexedObj.Interval,
			}
			errs.RelationshipsErrors = append(errs.RelationshipsErrors, msg)
			continue
		}

		authenticatedCtx := auth.InjectPrincipal(ctx, principal)

		handler := relationship.SetRelationshipHandler{}
		_, err := handler.Execute(authenticatedCtx, manager, &types.SetRelationshipRequest{
			PolicyId:     simCtx.Policy.Id,
			Relationship: indexedObj.Obj,
		})
		if err != nil {
			if errors.Is(err, errors.ErrorType_BAD_INPUT) {
				err := &types.LocatedMessage{
					Message:   err.Error(),
					Kind:      types.LocatedMessage_ERROR,
					InputName: "relationships",
					Interval:  indexedObj.Interval,
				}
				errs.RelationshipsErrors = append(errs.RelationshipsErrors, err)
			} else {
				return nil, err
			}
		}
	}
	return errs, nil
}

func HandleVerifyTheorem(ctx context.Context, manager runtime.RuntimeManager, req *types.VerifyTheoremsRequest) (*types.VerifyTheoremsResponse, error) {
	repository := NewSandboxRepository(manager.GetKVStore())

	record, err := repository.GetSandbox(ctx, req.Handle)
	if err != nil {
		return nil, newVerifyTheoremsErr(err, req.Handle)
	}
	if record == nil {
		return nil, newVerifyTheoremsErr(errors.Wrap("sandbox not found", errors.ErrorType_NOT_FOUND), req.Handle)
	}
	if !record.Initialized {
		return nil, newVerifyTheoremsErr(errors.Wrap("uninitialized sandbox cannot execute theorems", errors.ErrorType_OPERATION_FORBIDDEN), req.Handle)
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

	return &types.VerifyTheoremsResponse{
		Result: result,
	}, nil
}

func HandleGetCatalogue(ctx context.Context, manager runtime.RuntimeManager, req *types.GetCatalogueRequest) (*types.GetCatalogueResponse, error) {
	repository := NewSandboxRepository(manager.GetKVStore())

	record, err := repository.GetSandbox(ctx, req.Handle)
	if err != nil {
		return nil, newGetCatalogueErr(err, req.Handle)
	}
	if record == nil {
		return nil, newGetCatalogueErr(errors.Wrap("sandbox not found", errors.ErrorType_NOT_FOUND), req.Handle)
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

	return &types.GetCatalogueResponse{
		Catalogue: catalogue,
	}, nil
}

func HandleRestoreScratchpad(ctx context.Context, manager runtime.RuntimeManager, req *types.RestoreScratchpadRequest) (*types.RestoreScratchpadResponse, error) {
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

	return &types.RestoreScratchpadResponse{
		Scratchpad: record.Scratchpad,
	}, nil
}

func HandleGetSandbox(ctx context.Context, manager runtime.RuntimeManager, req *types.GetSandboxRequest) (*types.GetSandboxResponse, error) {
	repository := NewSandboxRepository(manager.GetKVStore())

	record, err := repository.GetSandbox(ctx, req.Handle)
	if err != nil {
		return nil, newGetSandboxErr(err, req.Handle)
	}
	if record == nil {
		err = errors.Wrap("sandbox not found", errors.ErrorType_NOT_FOUND)
		return nil, newGetSandboxErr(err, req.Handle)
	}
	return &types.GetSandboxResponse{
		Record: record,
	}, nil
}

// getPolicyEndPosition return a BufferPosition to the last character
// in pol. If an empty string, returns position 1,1
func getPolicyEndPosition(pol string) types.BufferPosition {
	pos := types.BufferPosition{
		Line:   1,
		Column: 1,
	}
	lines := strings.Split(pol, "\n")
	if len(lines) > 0 {
		lineCount := uint64(len(lines))
		lastLine := lines[lineCount-1]
		lastCol := uint64(len(lastLine))
		pos.Line = lineCount
		pos.Column = lastCol
	}
	return pos
}

func HandleGetSandboxSamples(ctx context.Context, _ runtime.RuntimeManager, req *types.GetSampleSandboxesRequest) (*types.GetSampleSandboxesResponse, error) {
	return &types.GetSampleSandboxesResponse{
		Samples: Samples,
	}, nil
}

func HandleExplainCheck(ctx context.Context, manager runtime.RuntimeManager, req *types.ExplainCheckRequest) (*types.ExplainCheckResponse, error) {
	repository := NewSandboxRepository(manager.GetKVStore())

	record, err := repository.GetSandbox(ctx, req.Handle)
	if err != nil {
		return nil, newExplainCheckError(err, req.Handle)
	}
	if record == nil {
		return nil, newExplainCheckError(errors.Wrap("sandbox not found", errors.ErrorType_NOT_FOUND), req.Handle)
	}
	if !record.Initialized {
		err := errors.Wrap("uninitialized sandbox cannot execute theorems",
			errors.ErrorType_OPERATION_FORBIDDEN)
		return nil, newExplainCheckError(err, req.Handle)
	}

	manager, err = GetManagerForSandbox(manager, req.Handle)
	if err != nil {
		return nil, newExplainCheckError(err, req.Handle)
	}

	op := &types.Operation{
		Object:     req.Object,
		Permission: req.Permission,
	}

	engine, err := zanzi.NewZanzi(manager.GetKVStore(), manager.GetLogger())
	if err != nil {
		return nil, newExplainCheckError(err, req.Handle)
	}
	authorized, tree, err := engine.ExplainCheck(ctx, record.Ctx.Policy, op, req.Actor)
	if err != nil {
		return nil, newExplainCheckError(err, req.Handle)
	}

	bytes, err := json.Marshal(tree)
	if err != nil {
		return nil, newExplainCheckError(err, req.Handle)
	}

	return &types.ExplainCheckResponse{
		Authorized: authorized,
		TreeJson:   bytes,
	}, nil
}

func HandleDOTExplainCheck(ctx context.Context, manager runtime.RuntimeManager, req *types.DOTExplainCheckRequest) (*types.DOTExplainCheckResponse, error) {
	repository := NewSandboxRepository(manager.GetKVStore())

	record, err := repository.GetSandbox(ctx, req.Handle)
	if err != nil {
		return nil, newExplainCheckError(err, req.Handle)
	}
	if record == nil {
		return nil, newExplainCheckError(errors.Wrap("sandbox not found", errors.ErrorType_NOT_FOUND), req.Handle)
	}
	if !record.Initialized {
		err := errors.Wrap("uninitialized sandbox cannot execute theorems",
			errors.ErrorType_OPERATION_FORBIDDEN)
		return nil, newExplainCheckError(err, req.Handle)
	}

	manager, err = GetManagerForSandbox(manager, req.Handle)
	if err != nil {
		return nil, newExplainCheckError(err, req.Handle)
	}

	op := &types.Operation{
		Object:     req.Object,
		Permission: req.Permission,
	}

	engine, err := zanzi.NewZanzi(manager.GetKVStore(), manager.GetLogger())
	if err != nil {
		return nil, newExplainCheckError(err, req.Handle)
	}
	authorized, tree, err := engine.DOTExplainCheck(ctx, record.Ctx.Policy, op, req.Actor)
	if err != nil {
		return nil, newExplainCheckError(err, req.Handle)
	}

	return &types.DOTExplainCheckResponse{
		Authorized: authorized,
		DotGraph:   tree,
	}, nil
}
