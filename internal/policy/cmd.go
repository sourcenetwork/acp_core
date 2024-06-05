package policy

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

type CreatePolicyHandler struct{}

// Execute consumes the data supplied in the command and creates a new ACP Policy and stores it in the given engine.
func (c *CreatePolicyHandler) Execute(ctx context.Context, runtime runtime.RuntimeManager, cmd *types.CreatePolicyRequest) (*types.CreatePolicyResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	principal, err := auth.ExtractPrincipal(ctx)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}
	if principal.IsAnonymous() {
		return nil, fmt.Errorf("CreatePolicy: requires authenticated principal: %w", auth.ErrUnauthenticatd)
	}
	identifier := principal.Identifier()

	ir, err := Unmarshal(cmd.Policy, cmd.MarshalType)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}

	err = basicPolicyIRSpec(&ir)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}

	counter := newPolicyCounter(runtime)
	i, err := counter.GetNextAndIncrement(ctx)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}

	factory := factory{}
	record, err := factory.Create(ir, identifier, i, cmd.CreationTime)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}

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
		Creator:    identifier,
		PolicyId:   record.Policy.Id,
		PolicyName: record.Policy.Name,
	}
	err = eventManager.EmitEvent(&event)
	if err != nil {
		return nil, err
	}

	return &types.CreatePolicyResponse{
		Policy: record.Policy,
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
