package policy

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/raccoon"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
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

	ir, err := Unmarshal(req.Policy, req.MarshalType)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}

	err = basicPolicyIRSpec(&ir)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
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

	factory := factory{}
	record, err := factory.Create(ir, req.Attributes, i, now)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}
	record.PolicyDefinition = req.Policy
	record.MarshalType = req.MarshalType

	spec := validPolicySpec{}
	err = spec.Satisfies(record.Policy)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
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
		Policy:     record.Policy,
		Attributes: record.Metadata,
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
