package services

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/authz_db"
	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/internal/relationship"
	"github.com/sourcenetwork/acp_core/internal/system"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ Decorator = (*MsgSpanDecorator)(nil)
var _ types.ACPEngineServer = (*acpEngine)(nil)

// acpEngine implements the ACP module MsgServer interface and accepts
// decorating functions which can wrap the execution of a Msg.
type acpEngine struct {
	hooks   []Decorator
	runtime runtime.RuntimeManager
}

// NewCmdSrever creates a message server for Embedded ACP
func NewACPEngine(runtime runtime.RuntimeManager, hooks ...Decorator) types.ACPEngineServer {
	return &acpEngine{
		hooks:   hooks,
		runtime: runtime,
	}
}

func (s *acpEngine) CreatePolicy(ctx context.Context, msg *types.CreatePolicyRequest) (*types.CreatePolicyResponse, error) {
	handler := policy.CreatePolicyHandler{}
	h := func(ctx context.Context, msg *types.CreatePolicyRequest) (*types.CreatePolicyResponse, error) {
		return handler.Execute(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, msg)
}

func (s *acpEngine) SetRelationship(ctx context.Context, msg *types.SetRelationshipRequest) (*types.SetRelationshipResponse, error) {
	handler := relationship.SetRelationshipHandler{}
	h := func(ctx context.Context, msg *types.SetRelationshipRequest) (*types.SetRelationshipResponse, error) {
		return handler.Execute(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, msg)
}

func (s *acpEngine) DeleteRelationship(ctx context.Context, msg *types.DeleteRelationshipRequest) (*types.DeleteRelationshipResponse, error) {
	handler := relationship.DeleteRelationshipHandler{}
	h := func(ctx context.Context, msg *types.DeleteRelationshipRequest) (*types.DeleteRelationshipResponse, error) {
		return handler.Execute(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, msg)
}

func (s *acpEngine) RegisterObject(ctx context.Context, msg *types.RegisterObjectRequest) (*types.RegisterObjectResponse, error) {
	handler := relationship.RegisterObjectHandler{}
	h := func(ctx context.Context, msg *types.RegisterObjectRequest) (*types.RegisterObjectResponse, error) {
		return handler.Execute(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, msg)
}

func (s *acpEngine) ArchiveObject(ctx context.Context, msg *types.ArchiveObjectRequest) (*types.ArchiveObjectResponse, error) {
	handler := relationship.ArchiveObjectHandler{}
	h := func(ctx context.Context, msg *types.ArchiveObjectRequest) (*types.ArchiveObjectResponse, error) {
		return handler.Execute(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, msg)
}

func (s *acpEngine) GetObjectRegistration(ctx context.Context, req *types.GetObjectRegistrationRequest) (*types.GetObjectRegistrationResponse, error) {
	h := func(ctx context.Context, msg *types.GetObjectRegistrationRequest) (*types.GetObjectRegistrationResponse, error) {
		return relationship.GetObjectRegistrationHandler(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, req)
}

func (s *acpEngine) FilterRelationships(ctx context.Context, req *types.FilterRelationshipsRequest) (*types.FilterRelationshipsResponse, error) {
	h := func(ctx context.Context, msg *types.FilterRelationshipsRequest) (*types.FilterRelationshipsResponse, error) {
		return relationship.FilterRelationshipsHandler(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, req)
}

func (s *acpEngine) DeletePolicy(ctx context.Context, req *types.DeletePolicyRequest) (*types.DeletePolicyResponse, error) {
	h := func(ctx context.Context, msg *types.DeletePolicyRequest) (*types.DeletePolicyResponse, error) {
		return policy.HandleDeletePolicy(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, req)
}

func (s *acpEngine) ValidatePolicy(ctx context.Context, req *types.ValidatePolicyRequest) (*types.ValidatePolicyResponse, error) {
	h := func(ctx context.Context, msg *types.ValidatePolicyRequest) (*types.ValidatePolicyResponse, error) {
		return policy.ValidatePolicy(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, req)
}

func (s *acpEngine) VerifyAccessRequest(ctx context.Context, req *types.VerifyAccessRequestRequest) (*types.VerifyAccessRequestResponse, error) {
	h := func(ctx context.Context, msg *types.VerifyAccessRequestRequest) (*types.VerifyAccessRequestResponse, error) {
		return authz_db.VerifyAccessRequest(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, req)
}

func (s *acpEngine) GetPolicy(ctx context.Context, req *types.GetPolicyRequest) (*types.GetPolicyResponse, error) {
	h := func(ctx context.Context, msg *types.GetPolicyRequest) (*types.GetPolicyResponse, error) {
		return policy.HandleGetPolicy(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, req)
}

func (s *acpEngine) ListPolicies(ctx context.Context, req *types.ListPoliciesRequest) (*types.ListPoliciesResponse, error) {
	h := func(ctx context.Context, msg *types.ListPoliciesRequest) (*types.ListPoliciesResponse, error) {
		return policy.ListPolicies(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, req)
}

func (s *acpEngine) TransferObject(ctx context.Context, req *types.TransferObjectRequest) (*types.TransferObjectResponse, error) {
	h := func(ctx context.Context, msg *types.TransferObjectRequest) (*types.TransferObjectResponse, error) {
		handler := relationship.TransferObjectHandler{}
		return handler.Execute(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, req)
}

func (s *acpEngine) SetParams(ctx context.Context, req *types.SetParamsRequest) (*types.SetParamsResponse, error) {
	h := func(ctx context.Context, msg *types.SetParamsRequest) (*types.SetParamsResponse, error) {
		return system.HandleSetParams(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, req)
}

func (s *acpEngine) GetParams(ctx context.Context, req *types.GetParamsRequest) (*types.GetParamsResponse, error) {
	h := func(ctx context.Context, msg *types.GetParamsRequest) (*types.GetParamsResponse, error) {
		return system.HandleGetParams(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, req)
}

func (s *acpEngine) EvaluateTheorem(ctx context.Context, req *types.EvaluateTheoremRequest) (*types.EvaluateTheoremResponse, error) {
	h := func(ctx context.Context, msg *types.EvaluateTheoremRequest) (*types.EvaluateTheoremResponse, error) {
		return policy.EvaluateTheorem(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, req)
}

func (s *acpEngine) GetPolicyCatalogue(ctx context.Context, req *types.GetPolicyCatalogueRequest) (*types.GetPolicyCatalogueResponse, error) {
	h := func(ctx context.Context, msg *types.GetPolicyCatalogueRequest) (*types.GetPolicyCatalogueResponse, error) {
		return policy.GetPolicyCatalogue(ctx, s.runtime, msg)
	}
	return applyMiddleware(ctx, h, s.hooks, req)
}

func (s *acpEngine) AmendRegistration(ctx context.Context, req *types.AmendRegistrationRequest) (*types.AmendRegistrationResponse, error) {
	h := func(ctx context.Context, req *types.AmendRegistrationRequest) (*types.AmendRegistrationResponse, error) {
		h := relationship.AmendRegistrationHandler{}
		return h.Handle(ctx, s.runtime, req)
	}
	return applyMiddleware(ctx, h, s.hooks, req)
}

func (s *acpEngine) UnarchiveObject(ctx context.Context, req *types.UnarchiveObjectRequest) (*types.UnarchiveObjectResponse, error) {
	h := func(ctx context.Context, req *types.UnarchiveObjectRequest) (*types.UnarchiveObjectResponse, error) {
		h := relationship.UnarchiveObjectHandler{}
		return h.Handle(ctx, s.runtime, req)
	}
	return applyMiddleware(ctx, h, s.hooks, req)
}
