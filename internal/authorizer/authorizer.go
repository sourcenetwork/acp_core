package authorizer

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// VerifyManagementPermission acts as an Authorization Request engine
// which validates whether a Relationship can be set or deleted by an Actor.
//
// The Relation evaluation is done through a CheckExpression call, which provides a request time
// expression. The expression is defined during acp_Core policy creation by a transformer.
//
// For instance, take the Relationship (obj:foo, reader, steve) being submitted by Actor Bob.
// Bob is allowed to Create that relationship if and only if:
// Bob has the permission _can_manage_reader for "obj:foo".
//
// A given Relationship is only valid if for the Relationship's Object and Relation
// the Actor has an associated permission to manage the Object, Relation pair.
func VerifyManagementPermission(ctx context.Context, engine *zanzi.Adapter, policy *types.Policy, obj *types.Object, relation string, actor *types.Actor) (bool, error) {
	resource := policy.GetResourceByName(obj.Resource)
	if resource == nil {
		return false, errors.New("resource not found in policy", errors.ErrorType_NOT_FOUND,
			errors.Pair("policy", policy.Id),
			errors.Pair("resource", obj.Resource),
		)
	}
	perm := resource.GetManagementPermissionByName(relation)
	if perm == nil {
		return false, errors.New("management permission not found for relation in resource", errors.ErrorType_NOT_FOUND,
			errors.Pair("policy", policy.Id),
			errors.Pair("resource", obj.Resource),
			errors.Pair("relation", relation),
		)
	}

	return engine.CheckExpression(ctx, policy, obj, perm.Expression, actor)
}

// CheckManagementAuthority implements a handler for CheckManagementAuthority requests
func CheckManagementAuthority(ctx context.Context, runtime runtime.RuntimeManager, req *types.CheckManagementAuthorityRequest) (*types.CheckManagementAuthorityResponse, error) {
	zanzi, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	if req.Actor == nil {
		return nil, errors.New("actor cannot be nil", errors.ErrorType_BAD_INPUT)
	}
	if req.Object == nil {
		return nil, errors.New("object cannot be nil", errors.ErrorType_BAD_INPUT)
	}
	if req.Relation == "" {
		return nil, errors.New("relation cannot be empty", errors.ErrorType_BAD_INPUT)
	}

	rec, err := zanzi.GetPolicy(ctx, req.PolicyId)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, errors.Wrap("checking management authority", errors.ErrPolicyNotFound(req.PolicyId))
	}

	authorized, err := VerifyManagementPermission(ctx, zanzi, rec.Policy, req.Object, req.Relation, req.Actor)
	if err != nil {
		return nil, errors.Wrap("checking management authority", err)
	}

	return &types.CheckManagementAuthorityResponse{
		Authorized: authorized,
	}, nil
}
