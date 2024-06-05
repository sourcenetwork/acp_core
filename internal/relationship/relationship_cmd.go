package relationship

import (
	"context"
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/did"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

type SetRelationshipHandler struct{}

func (c *SetRelationshipHandler) Execute(ctx context.Context, runtime runtime.RuntimeManager, cmd *types.SetRelationshipRequest) (*types.SetRelationshipResponse, error) {
	//eventManager := runtime.GetEventManager()

	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	principal, err := auth.ExtractPrincipalWithType(ctx, auth.DID)
	if err != nil {
		return nil, fmt.Errorf("MsgRegisterObject: %w", err)
	}
	did := principal.Identifier()

	authorizer := NewRelationshipAuthorizer(engine)

	rec, err := engine.GetPolicy(ctx, cmd.PolicyId)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, fmt.Errorf("MsgSetRelationship: policy %v: %w", cmd.PolicyId, types.ErrPolicyNotFound)
	}
	policy := rec.Policy

	err = c.validate(ctx, policy, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to set relationship: %w", err)
	}

	creatorActor := types.Actor{
		Id: did,
	}

	obj := cmd.Relationship.Object
	ownerRecord, err := QueryOwnerRelationship(ctx, engine, policy, obj)
	if err != nil {
		return nil, fmt.Errorf("failed to set relationship: %w", err)
	}
	if ownerRecord == nil {
		return nil, fmt.Errorf("failed to set relationship: object %v: %w", obj, types.ErrObjectNotFound)
	}

	authorized, err := authorizer.IsAuthorized(ctx, policy, cmd.Relationship, &creatorActor)
	if err != nil {
		return nil, fmt.Errorf("failed to set relationship: %w", err)
	}
	if !authorized {
		return nil, fmt.Errorf("failed to set relationship: actor %v is not a manager of relation %v for object %v: %w",
			did, cmd.Relationship.Relation, obj, types.ErrNotAuthorized)
	}

	record, err := engine.GetRelationship(ctx, policy, cmd.Relationship)
	if err != nil {
		return nil, fmt.Errorf("failed to set relationship: %w", err)
	}
	if record != nil {
		return &types.SetRelationshipResponse{
			RecordExisted: true,
			Record:        record,
		}, nil
	}

	record = &types.RelationshipRecord{
		PolicyId:     policy.Id,
		Relationship: cmd.Relationship,
		CreationTime: cmd.CreationTime,
		Actor:        did,
		Archived:     false,
	}
	_, err = engine.SetRelationship(ctx, policy, record)
	if err != nil {
		return nil, fmt.Errorf("failed to set relationship: %w", err)
	}

	return &types.SetRelationshipResponse{
		RecordExisted: false,
		Record:        record,
	}, nil
}

func (c *SetRelationshipHandler) validate(ctx context.Context, pol *types.Policy, cmd *types.SetRelationshipRequest) error {
	err := relationshipSpec(pol, cmd.Relationship)
	if err != nil {
		return err
	}

	if cmd.Relationship.Relation == policy.OwnerRelation {
		return ErrSetOwnerRel
	}

	if actor := cmd.Relationship.Subject.GetActor(); actor != nil {
		if err := did.IsValidDID(actor.Id); err != nil {
			return err
		}
	}

	return nil
}

type DeleteRelationshipHandler struct{}

func (c *DeleteRelationshipHandler) Execute(ctx context.Context, runtime runtime.RuntimeManager, cmd *types.DeleteRelationshipRequest) (*types.DeleteRelationshipResponse, error) {
	//eventManager := runtime.GetEventManager()

	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	principal, err := auth.ExtractPrincipalWithType(ctx, auth.DID)
	if err != nil {
		return nil, fmt.Errorf("MsgRegisterObject: %w", err)
	}
	did := principal.Identifier()

	authorizer := NewRelationshipAuthorizer(engine)

	err = c.validate(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to delete relationship: %w", err)
	}

	rec, err := engine.GetPolicy(ctx, cmd.PolicyId)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, fmt.Errorf("MsgDeleteRelationship: policy %v: %w", cmd.PolicyId, types.ErrPolicyNotFound)
	}
	policy := rec.Policy

	isAuthorized, err := c.isActorAuthorized(ctx, authorizer, policy, cmd, did)
	if err != nil {
		return nil, fmt.Errorf("failed to delete relationship: %w", err)
	}

	if !isAuthorized {
		return nil, fmt.Errorf("failed to delete relationship: %w", types.ErrNotAuthorized)
	}

	found, err := engine.DeleteRelationship(ctx, policy, cmd.Relationship)
	if err != nil {
		return nil, fmt.Errorf("failed to delete relationship: %w", err)
	}

	return &types.DeleteRelationshipResponse{
		RecordFound: bool(found),
	}, nil

}

func (c *DeleteRelationshipHandler) validate(ctx context.Context, cmd *types.DeleteRelationshipRequest) error {
	if cmd.Relationship == nil {
		return types.ErrRelationshipNil
	}

	if cmd.Relationship.Relation == policy.OwnerRelation {
		return ErrDeleteOwnerRel
	}

	return nil
}

// verifies whether actor is authorized to remove the specified Relationship
func (c *DeleteRelationshipHandler) isActorAuthorized(ctx context.Context, authorizer *RelationshipAuthorizer, policy *types.Policy, cmd *types.DeleteRelationshipRequest, initiator string) (bool, error) {
	creatorActor := types.Actor{
		Id: initiator,
	}
	return authorizer.IsAuthorized(ctx, policy, cmd.Relationship, &creatorActor)
}
