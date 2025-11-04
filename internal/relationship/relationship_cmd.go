package relationship

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/authorizer"
	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

type SetRelationshipHandler struct{}

func (c *SetRelationshipHandler) Execute(ctx context.Context, runtime runtime.RuntimeManager, cmd *types.SetRelationshipRequest) (*types.SetRelationshipResponse, error) {
	//eventManager := runtime.GetEventManager()
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, newSetRelationshipErr(err)
	}

	principal, err := auth.ExtractPrincipalWithType(ctx, types.PrincipalKind_DID)
	if err != nil {
		return nil, newSetRelationshipErr(err)
	}
	did := principal.Identifier

	rec, err := engine.GetPolicy(ctx, cmd.PolicyId)
	if err != nil {
		return nil, newSetRelationshipErr(err)
	}
	if rec == nil {
		return nil, newSetRelationshipErr(errors.ErrPolicyNotFound(cmd.PolicyId))
	}
	policy := rec.Policy

	err = c.validate(policy, cmd)
	if err != nil {
		return nil, newSetRelationshipErr(err)
	}

	creatorActor := types.Actor{
		Id: did,
	}

	obj := cmd.Relationship.Object
	ownerRecord, err := queryOwnerRelationship(ctx, engine, policy, obj)
	if err != nil {
		return nil, newSetRelationshipErr(err)
	}
	if ownerRecord == nil {
		return nil, newSetRelationshipErr(errors.New("cannot set relationship for unregistered object", errors.ErrorType_NOT_FOUND,
			errors.Pair("policy", cmd.PolicyId),
			errors.Pair("resource", cmd.Relationship.Object.Resource),
			errors.Pair("object", cmd.Relationship.Object.Id),
		))
	}

	authorized, err := authorizer.VerifyManagementPermission(ctx, engine, policy,
		cmd.Relationship.Object, cmd.Relationship.Relation, &creatorActor)
	if err != nil {
		return nil, newSetRelationshipErr(err)
	}
	if !authorized {
		return nil, newSetRelationshipErr(
			errors.New("cannot create relationship: actor is not a manager of relation", errors.ErrorType_UNAUTHORIZED,
				errors.Pair("policy", cmd.PolicyId),
				errors.Pair("relation", cmd.Relationship.Relation),
				errors.Pair("actor", did),
			))
	}
	record, err := engine.GetRelationship(ctx, policy, cmd.Relationship)
	if err != nil {
		return nil, newSetRelationshipErr(err)
	}
	if record != nil {
		return &types.SetRelationshipResponse{
			RecordExisted: true,
			Record:        record,
		}, nil
	}

	ts, err := runtime.GetTimeService().GetNow(ctx)
	if err != nil {
		return nil, newSetRelationshipErr(err)
	}
	record = &types.RelationshipRecord{
		PolicyId:     policy.Id,
		Relationship: cmd.Relationship,
		Archived:     false,
		Metadata: &types.RecordMetadata{
			Creator:    &principal,
			CreationTs: ts,
			Supplied:   cmd.Metadata,
		},
	}
	_, err = engine.SetRelationship(ctx, policy, record)
	if err != nil {
		return nil, newSetRelationshipErr(err)
	}

	return &types.SetRelationshipResponse{
		RecordExisted: false,
		Record:        record,
	}, nil
}

func (c *SetRelationshipHandler) validate(pol *types.Policy, cmd *types.SetRelationshipRequest) error {
	err := relationshipSpec(pol, cmd.Relationship)
	if err != nil {
		return err
	}
	if cmd.Relationship.Relation == policy.OwnerRelation {
		return ErrSetOwnerRel
	}
	return nil
}

type DeleteRelationshipHandler struct{}

func (c *DeleteRelationshipHandler) Execute(ctx context.Context, runtime runtime.RuntimeManager, cmd *types.DeleteRelationshipRequest) (*types.DeleteRelationshipResponse, error) {
	//eventManager := runtime.GetEventManager()
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, newDeleteRelationshipErr(err)
	}

	principal, err := auth.ExtractPrincipalWithType(ctx, types.PrincipalKind_DID)
	if err != nil {
		return nil, newDeleteRelationshipErr(err)
	}
	did := principal.Identifier

	err = c.validate(cmd)
	if err != nil {
		return nil, newDeleteRelationshipErr(err)
	}

	rec, err := engine.GetPolicy(ctx, cmd.PolicyId)
	if err != nil {
		return nil, newDeleteRelationshipErr(err)
	}
	if rec == nil {
		return nil, newDeleteRelationshipErr(errors.ErrPolicyNotFound(cmd.PolicyId))
	}
	policy := rec.Policy

	isAuthorized, err := authorizer.VerifyManagementPermission(ctx, engine, policy,
		cmd.Relationship.Object, cmd.Relationship.Relation, types.NewActor(did))
	if err != nil {
		return nil, newDeleteRelationshipErr(err)
	}

	if !isAuthorized {
		return nil, newDeleteRelationshipErr(errors.Wrap("cannot delete relationship: actor is not a manager of relation", errors.ErrorType_UNAUTHORIZED,
			errors.Pair("policy", cmd.PolicyId),
			errors.Pair("relation", cmd.Relationship.Relation),
			errors.Pair("actor", did),
		))
	}

	found, err := engine.DeleteRelationship(ctx, policy, cmd.Relationship)
	if err != nil {
		return nil, newDeleteRelationshipErr(err)
	}

	return &types.DeleteRelationshipResponse{
		RecordFound: bool(found),
	}, nil
}

func (c *DeleteRelationshipHandler) validate(cmd *types.DeleteRelationshipRequest) error {
	if cmd.Relationship.Relation == policy.OwnerRelation {
		return ErrDeleteOwnerRel
	}
	return nil
}
