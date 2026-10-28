package services

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/alloy"
	"github.com/sourcenetwork/acp_core/internal/decorator"
	"github.com/sourcenetwork/acp_core/internal/policy_compiler"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ types.PolicyCompilerServer = (*PolicyCompilerService)(nil)

// PolicyCompilerService compiles policies into Alloy models
type PolicyCompilerService struct {
	decorator decorator.Decorator
}

func NewPolicyCompilerService() *PolicyCompilerService {
	return &PolicyCompilerService{
		decorator: decorator.RecoverDecorator,
	}
}

// Compile processes a policy and emits it as a bundle of Alloy modules.
// An unspecified marshal type defaults to YAML and an empty module name to the policy name.
// If the request has assumptions, they are emitted as an additional module.
func (s *PolicyCompilerService) Compile(ctx context.Context, req *types.CompileRequest) (*types.CompileResponse, error) {
	h := func(ctx context.Context, req *types.CompileRequest) (*types.CompileResponse, error) {
		marshalType := req.MarshalType
		if marshalType == types.PolicyMarshalingType_UNKNOWN {
			marshalType = types.PolicyMarshalingType_YAML
		}
		model, err := policy_compiler.Compile(req.Policy, marshalType)
		if err != nil {
			return nil, err
		}
		files := alloy.Emit(model, req.Module)
		if len(req.Constraints) > 0 {
			assumptions, err := policy_compiler.LowerAssumptions(model, req.Constraints)
			if err != nil {
				return nil, err
			}
			files = append(files, alloy.EmitAssumptions(model, req.Module, assumptions))
		}

		resp := &types.CompileResponse{}
		for _, f := range files {
			resp.Files = append(resp.Files, &types.CompiledFile{Name: f.Name, Content: f.Content})
		}
		return resp, nil
	}
	return decorator.DecorateTypedHandler(h, s.decorator)(ctx, req)
}
