package playground2

import (
	"context"
	"fmt"
	"strings"

	"github.com/sourcenetwork/acp_core/internal/sandbox"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var ErrScenario = errors.New("scenario", errors.ErrorType_BAD_INPUT)

// emptyAssertions is loaded alongside a world, which carries no assertions to evaluate here
const emptyAssertions = "Authorizations {}\nDelegations {}"

// Check returns whether an actor has a permission or relation on an object, in a world.
// ctx must carry request data (see telemetry.InitRequestContext).
func Check(ctx context.Context, req *types.CheckRequest) (*types.CheckResponse, error) {
	engine, policy, err := loadScenario(ctx, req.Scenario)
	if err != nil {
		return nil, err
	}
	op := &types.Operation{Object: req.Object, Permission: req.Permission}
	allowed, err := engine.Check(ctx, policy, op, req.Actor)
	if err != nil {
		return nil, err
	}
	return &types.CheckResponse{Allowed: allowed}, nil
}

// Explain performs a Check, returning the engine's search graph alongside the result.
// ctx must carry request data (see telemetry.InitRequestContext).
func Explain(ctx context.Context, req *types.ExplainRequest) (*types.ExplainResponse, error) {
	engine, policy, err := loadScenario(ctx, req.Scenario)
	if err != nil {
		return nil, err
	}
	op := &types.Operation{Object: req.Object, Permission: req.Permission}
	allowed, graph, err := engine.ExplainCheck(ctx, policy, op, req.Actor)
	if err != nil {
		return nil, err
	}
	return &types.ExplainResponse{Allowed: allowed, Graph: graph}, nil
}

// loadScenario loads a world into an ephemeral sandbox, returning an engine over it and the processed policy.
// System constraints play no part: the engine evaluates whatever relationships the world has.
func loadScenario(ctx context.Context, scenario *types.Scenario) (*zanzi.Adapter, *types.Policy, error) {
	if scenario == nil || scenario.SystemModel == nil {
		return nil, nil, errors.Wrap("system model is required", ErrScenario)
	}
	var relationships string
	if scenario.World != nil {
		relationships = scenario.World.Relationships
	}

	manager, err := runtime.NewRuntimeManager()
	if err != nil {
		return nil, nil, err
	}
	newResp, err := sandbox.HandleNewSandboxRequest(ctx, manager, &types.NewSandboxRequest{})
	if err != nil {
		return nil, nil, err
	}
	handler := sandbox.SetStateHandler{}
	setResp, err := handler.Handle(ctx, manager, &types.SetStateRequest{
		Handle: newResp.Record.Handle,
		Data: &types.SandboxData{
			PolicyDefinition: scenario.SystemModel.Policy,
			Relationships:    relationships,
			PolicyTheorem:    emptyAssertions,
		},
	})
	if err != nil {
		return nil, nil, err
	}
	if setResp.Errors.HasErrors() {
		return nil, nil, errors.Wrap(formatLoadErrors(setResp.Errors), ErrScenario)
	}

	sandboxManager, err := sandbox.GetManagerForSandbox(manager, newResp.Record.Handle)
	if err != nil {
		return nil, nil, err
	}
	engine, err := zanzi.NewZanzi(sandboxManager.GetKVStore(), sandboxManager.GetLogger())
	if err != nil {
		return nil, nil, err
	}
	return engine, setResp.Record.Ctx.Policy, nil
}

// formatLoadErrors renders the errors found while loading a scenario, eg. `relationships 2:0: object not registered`
func formatLoadErrors(errs *types.SandboxDataErrors) string {
	var msgs []string
	add := func(input string, located []*types.LocatedMessage) {
		for _, msg := range located {
			position := ""
			if msg.Interval != nil && msg.Interval.Start != nil {
				position = fmt.Sprintf(" %v:%v", msg.Interval.Start.Line, msg.Interval.Start.Column)
			}
			msgs = append(msgs, fmt.Sprintf("%v%v: %v", input, position, msg.Message))
		}
	}
	add("policy", errs.PolicyErrors)
	add("relationships", errs.RelationshipsErrors)
	return "invalid scenario: " + strings.Join(msgs, "; ")
}
