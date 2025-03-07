package authorizer

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// ManagementRequest models an Actor's request to modify
// a set of Relationships for a given object and relation
type ManagementRequest struct {
	Policy   *types.Policy
	Object   *types.Object
	Relation string
	Actor    *types.Actor
}

func NewOperationAuthorizer(engine *zanzi.Adapter) *OperationAuthorizer {
	return &OperationAuthorizer{
		engine: engine,
	}
}

// OperationAuthorizer acts as an Authorization Request engine
// which validates whether a Relationship can be set or deleted by an Actor.
//
// The Permission evaluation is done through a Check call using the auxiliary permissions
// auto generated by the ACP module and attached to a permission.
//
// For instance, take the Relationship (obj:foo, reader, steve) being submitted by Actor Bob.
// Bob is allowed to Create that relationship if and only if:
// Bob has the permission _can_manage_reader for "obj:foo".
type OperationAuthorizer struct {
	engine *zanzi.Adapter
}

// IsAuthorized validates the given management request
//
// A given Relationship is only valid if for the Relationship's Object and Relation
// the Actor has an associated permission to manage the Object, Relation pair.
func (a *OperationAuthorizer) IsAuthorized(ctx context.Context, request *ManagementRequest) (bool, error) {
	resource := request.Policy.GetResourceByName(request.Object.Resource)
	if resource == nil {
		return false, errors.New("resource not found in policy", errors.ErrorType_NOT_FOUND,
			errors.Pair("policy", request.Policy.Id),
			errors.Pair("resource", request.Object.Resource),
		)
	}
	relation := resource.GetRelationByName(request.Relation)
	if relation == nil {
		return false, errors.New("relation not found in resource", errors.ErrorType_NOT_FOUND,
			errors.Pair("policy", request.Policy.Id),
			errors.Pair("resource", request.Object.Resource),
			errors.Pair("relation", request.Relation),
		)
	}

	authRequest := &types.Operation{
		Object:     request.Object,
		Permission: request.Policy.GetManagementPermissionName(request.Relation),
	}

	return a.engine.Check(ctx, request.Policy, authRequest, request.Actor)
}
