package policy

import (
	"context"
	"fmt"

	"github.com/cosmos/gogoproto/proto"
	"github.com/sourcenetwork/acp_core/internal/policy/ppp"
	"github.com/sourcenetwork/acp_core/internal/raccoon"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

type CreatePolicyHandler struct{}

// Execute consumes the data supplied in the command and creates a new ACP Policy and stores it in the given engine.
func (c *CreatePolicyHandler) Execute(ctx context.Context, runtime runtime.RuntimeManager, req *types.CreatePolicyRequest) (*types.CreatePolicyResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	counter := raccoon.NewCounterStoreFromRuntimeManager(runtime, policyCounterPrefix)
	releaser := counter.Acquire()
	defer releaser.Release()
	i, err := counter.GetNextAndIncrement(ctx)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}

	now, err := runtime.GetTimeService().GetNow(ctx)
	if err != nil {
		return nil, err
	}

	principal, err := auth.ExtractAuthenticatedPrincipal(ctx)
	if err != nil {
		return nil, err
	}

	policy, err := Unmarshal(req.Policy, req.MarshalType)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}

	pipeline := ppp.CreatePolicyPipelineFactory(i, policy.SpecificationType)
	result, err := pipeline.Process(policy)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}
	policy = &result.Policy

	record := &types.PolicyRecord{
		Policy:           policy,
		PolicyDefinition: req.Policy,
		MarshalType:      req.MarshalType,
		Metadata: &types.RecordMetadata{
			Creator:    &principal,
			CreationTs: now,
			Supplied:   req.Metadata,
		},
	}

	err = engine.SetPolicy(ctx, record)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}

	eventManager := runtime.GetEventManager()
	event := types.EventPolicyCreated{
		PolicyId:   record.Policy.Id,
		PolicyName: record.Policy.Name,
	}
	err = eventManager.EmitEvent(&event)
	if err != nil {
		return nil, err
	}

	return &types.CreatePolicyResponse{
		Record: record,
		Log:    result.Messages,
	}, nil
}

func HandleDeletePolicy(ctx context.Context, runtime runtime.RuntimeManager, req *types.DeletePolicyRequest) (*types.DeletePolicyResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, fmt.Errorf("delete policy: %w", err)
	}

	found, err := engine.DeletePolicy(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("delete policy: %v", err)
	}

	return &types.DeletePolicyResponse{
		Found: found,
	}, nil
}

type CreatePolicyWithSpecHandler struct{}

// Execute consumes the data supplied in the command and creates a new ACP Policy and stores it in the given engine.
func (c *CreatePolicyWithSpecHandler) Execute(ctx context.Context, runtime runtime.RuntimeManager, req *types.CreatePolicyWithSpecificationRequest) (*types.CreatePolicyWithSpecificationResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	counter := raccoon.NewCounterStoreFromRuntimeManager(runtime, policyCounterPrefix)
	releaser := counter.Acquire()
	defer releaser.Release()
	i, err := counter.GetNextAndIncrement(ctx)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}

	now, err := runtime.GetTimeService().GetNow(ctx)
	if err != nil {
		return nil, err
	}

	principal, err := auth.ExtractPrincipal(ctx)
	if err != nil {
		return nil, err
	}

	policy, err := Unmarshal(req.Policy, req.MarshalType)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}

	// TODO if create strategy is reject then error
	if policy.SpecificationType == types.PolicySpecificationType_NO_SPEC {
		policy.SpecificationType = req.RequiredSpec
	}

	if policy.SpecificationType != req.RequiredSpec {
		return nil, errors.Wrap("CreatePolicy: invalid specification type", errors.ErrorType_BAD_INPUT,
			errors.Pair("expected", req.RequiredSpec.String()),
			errors.Pair("got", policy.SpecificationType.String()),
		)
	}

	pipeline := ppp.CreatePolicyPipelineFactory(i, req.RequiredSpec)
	result, err := pipeline.Process(policy)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}
	policy = &result.Policy

	record := &types.PolicyRecord{
		Policy:           policy,
		PolicyDefinition: req.Policy,
		MarshalType:      req.MarshalType,
		Metadata: &types.RecordMetadata{
			Creator:    &principal,
			CreationTs: now,
			Supplied:   req.Metadata,
		},
	}

	err = engine.SetPolicy(ctx, record)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}

	eventManager := runtime.GetEventManager()
	event := types.EventPolicyCreated{
		PolicyId:   record.Policy.Id,
		PolicyName: record.Policy.Name,
	}
	err = eventManager.EmitEvent(&event)
	if err != nil {
		return nil, err
	}

	return &types.CreatePolicyWithSpecificationResponse{
		Record: record,
	}, nil
}

type EditPolicyHandler struct{}

func (h *EditPolicyHandler) Execute(ctx context.Context, runtime runtime.RuntimeManager, req *types.EditPolicyRequest) (*types.EditPolicyResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	oldRecord, err := engine.GetPolicy(ctx, req.PolicyId)
	if err != nil {
		return nil, fmt.Errorf("edit policy: %w", err)
	}
	if oldRecord == nil {
		return nil, errors.ErrPolicyNotFound(req.PolicyId)
	}

	principal, err := auth.ExtractAuthenticatedPrincipal(ctx)
	if err != nil {
		return nil, err
	}

	if !principal.Equals(oldRecord.Metadata.Creator) {
		return nil, errors.New("only the policy creator may edit it", errors.ErrorType_UNAUTHORIZED)
	}

	policy, err := Unmarshal(req.Policy, req.MarshalType)
	if err != nil {
		return nil, fmt.Errorf("EditPolicy: %w", err)
	}
	// TODO if create strategy is reject then error
	if policy.SpecificationType == types.PolicySpecificationType_NO_SPEC {
		policy.SpecificationType = types.PolicySpecificationType_NO_SPEC
	}

	pipeline := ppp.EditPolicyPipelineFactory(oldRecord.Policy)
	result, err := pipeline.Process(policy)
	if err != nil {
		return nil, fmt.Errorf("EditPolicy: %w", err)
	}
	policy = &result.Policy

	now, err := runtime.GetTimeService().GetNow(ctx)
	if err != nil {
		return nil, err
	}
	record := proto.Clone(oldRecord).(*types.PolicyRecord)
	record.Policy = policy
	record.MarshalType = req.MarshalType
	record.Metadata.LastModified = now
	record.PolicyDefinition = req.Policy

	count, err := engine.EditPolicy(ctx, record)
	if err != nil {
		return nil, fmt.Errorf("EditPolicy: %w", err)
	}

	return &types.EditPolicyResponse{
		RelatinshipsRemoved: count,
		Record:              record,
	}, nil
}

type EditPolicyMetadataHandler struct{}

func (h *EditPolicyMetadataHandler) Execute(ctx context.Context, runtime runtime.RuntimeManager, req *types.EditPolicyMetadataRequest) (*types.EditPolicyMetadataResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())

	oldRecord, err := engine.GetPolicy(ctx, req.PolicyId)
	if err != nil {
		return nil, fmt.Errorf("edit policy metadata: %w", err)
	}
	if oldRecord == nil {
		return nil, errors.ErrPolicyNotFound(req.PolicyId)
	}

	principal, err := auth.ExtractAuthenticatedPrincipal(ctx)
	if err != nil {
		return nil, err
	}

	if !principal.Equals(oldRecord.Metadata.Creator) {
		return nil, errors.New("only the policy creator may edit it", errors.ErrorType_UNAUTHORIZED)
	}

	now, err := runtime.GetTimeService().GetNow(ctx)
	if err != nil {
		return nil, err
	}
	record := proto.Clone(oldRecord).(*types.PolicyRecord)
	record.Metadata.LastModified = now
	record.Metadata.Supplied = req.Metadata

	_, err = engine.EditPolicy(ctx, record)
	if err != nil {
		return nil, fmt.Errorf("EditPolicyMetadata: %w", err)
	}

	return &types.EditPolicyMetadataResponse{
		Record: record,
	}, nil
}
