package relationship

import (
	"context"
	"fmt"

	prototypes "github.com/cosmos/gogoproto/types"

	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// RegisterObjectHandler creates an "owner" Relationship for the given object and subject,
// if the object does not have a previous owner.
// If the relationship exists but is archived by the same actor, unarchives it
// if relationship is active this command is a noop
type RegisterObjectHandler struct{}

func (c *RegisterObjectHandler) Execute(ctx context.Context, runtime runtime.RuntimeManager, cmd *types.RegisterObjectRequest) (*types.RegisterObjectResponse, error) {
	eventManager := runtime.GetEventManager()

	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	principal, err := auth.ExtractPrincipalWithType(ctx, auth.DID)
	if err != nil {
		return nil, fmt.Errorf("MsgRegisterObject: %w", err)
	}
	did := principal.Identifier()

	registration := &types.Registration{
		Object: cmd.Object,
		Actor: &types.Actor{
			Id: principal.Identifier(),
		},
	}

	rec, err := engine.GetPolicy(ctx, cmd.PolicyId)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, fmt.Errorf("MsgRegisterObject: policy %v: %w", cmd.PolicyId, types.ErrPolicyNotFound)
	}
	policy := rec.Policy

	err = c.validate(ctx, policy, did, registration)
	if err != nil {
		return nil, fmt.Errorf("failed to register object: %w", err)
	}

	err = nil
	var result types.RegistrationResult

	record, err := QueryOwnerRelationship(ctx, engine, policy, registration.Object)
	if err != nil {
		return nil, fmt.Errorf("failed to register object: %w", err)
	}

	switch c.resolveObjectStatus(record) {
	case statusUnregistered:
		result, err = c.unregisteredStrategy(ctx, engine, policy, registration, did, cmd.CreationTime)
	case statusArchived:
		result, err = c.archivedObjectStrategy(ctx, engine, policy, record, registration)
	case statusActive:
		result, err = c.activeObjectStrategy(record, registration)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to register object: %w", err)
	}

	record, err = QueryOwnerRelationship(ctx, engine, policy, registration.Object)
	if err != nil {
		return nil, fmt.Errorf("failed to register object: %w", err)
	}

	// TODO efactor the event type
	if result == types.RegistrationResult_Registered {
		eventManager.EmitEvent(&types.EventObjectRegistered{
			Actor:          registration.Actor.Id,
			PolicyId:       policy.Id,
			ObjectResource: registration.Object.Resource,
			ObjectId:       registration.Object.Id,
		})
	} else if result == types.RegistrationResult_Unarchived {
		eventManager.EmitEvent(&types.EventObjectUnarchived{
			Actor:          registration.Actor.Id,
			PolicyId:       policy.Id,
			ObjectResource: registration.Object.Resource,
			ObjectId:       registration.Object.Id,
		})
	}

	return &types.RegisterObjectResponse{
		Result: result,
		Record: record,
	}, nil
}

// validates the command input params
func (c *RegisterObjectHandler) validate(ctx context.Context, policy *types.Policy, creator string, registration *types.Registration) error {
	if policy == nil {
		return types.ErrPolicyNil
	}

	if err := registrationSpec(registration); err != nil {
		return err
	}

	return nil
}

func (c *RegisterObjectHandler) resolveObjectStatus(record *types.RelationshipRecord) objectRegistrationStatus {
	if record == nil {
		return statusUnregistered
	}
	if record.Archived {
		return statusArchived
	}
	return statusActive
}

func (c *RegisterObjectHandler) unregisteredStrategy(ctx context.Context, zanzi *zanzi.Adapter, pol *types.Policy, registration *types.Registration, creator string, creationTs *prototypes.Timestamp) (types.RegistrationResult, error) {
	err := c.createOwnerRelationship(ctx, zanzi, pol, registration, creator, creationTs)
	if err != nil {
		return types.RegistrationResult_NoOp, err
	}

	return types.RegistrationResult_Registered, nil
}

func (c *RegisterObjectHandler) createOwnerRelationship(ctx context.Context, zanzi *zanzi.Adapter, pol *types.Policy, registration *types.Registration, creator string, creationTs *prototypes.Timestamp) error {
	record := types.RelationshipRecord{
		Relationship: &types.Relationship{
			Object:   registration.Object,
			Relation: policy.OwnerRelation,
			Subject: &types.Subject{
				Subject: &types.Subject_Actor{
					Actor: registration.Actor,
				},
			},
		},
		Actor:        registration.Actor.Id,
		PolicyId:     pol.Id,
		Archived:     false,
		CreationTime: creationTs,
	}
	_, err := zanzi.SetRelationship(ctx, pol, &record)
	return err
}

func (c *RegisterObjectHandler) activeObjectStrategy(record *types.RelationshipRecord, registration *types.Registration) (types.RegistrationResult, error) {
	if record.Actor != registration.Actor.Id {
		return types.RegistrationResult_NoOp, types.ErrNotAuthorized
	}

	return types.RegistrationResult_NoOp, nil
}

func (c *RegisterObjectHandler) archivedObjectStrategy(ctx context.Context, engine *zanzi.Adapter, policy *types.Policy, record *types.RelationshipRecord, registration *types.Registration) (types.RegistrationResult, error) {
	if record.Actor != registration.Actor.Id {
		return types.RegistrationResult_NoOp, types.ErrNotAuthorized
	}

	err := c.unarchiveRelationship(ctx, engine, policy, record)
	if err != nil {
		return types.RegistrationResult_NoOp, err
	}

	return types.RegistrationResult_Unarchived, nil
}

func (c *RegisterObjectHandler) unarchiveRelationship(ctx context.Context, engine *zanzi.Adapter, policy *types.Policy, record *types.RelationshipRecord) error {
	record.Archived = false
	_, err := engine.SetRelationship(ctx, policy, record)
	return err
}

type UnregisterObjectHandler struct{}

func (c *UnregisterObjectHandler) Execute(ctx context.Context, runtime runtime.RuntimeManager, cmd *types.UnregisterObjectRequest) (*types.UnregisterObjectResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, err
	}

	principal, err := auth.ExtractPrincipalWithType(ctx, auth.DID)
	if err != nil {
		return nil, fmt.Errorf("MsgRegisterObject: %w", err)
	}
	did := principal.Identifier()

	rec, err := engine.GetPolicy(ctx, cmd.PolicyId)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, fmt.Errorf("MsgUnregisterObject: policy %v: %w", cmd.PolicyId, types.ErrPolicyNotFound)
	}
	pol := rec.Policy

	// TODO return found or not
	err = c.validate(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to unregister object: %w", err)
	}

	authorizer := NewRelationshipAuthorizer(engine)

	authRelationship := types.NewActorRelationship(cmd.Object.Resource, cmd.Object.Id, policy.OwnerRelation, did)
	actor := types.Actor{Id: did}
	authorized, err := authorizer.IsAuthorized(ctx, pol, authRelationship, &actor)
	if err != nil {
		return nil, fmt.Errorf("failed to unregister object: %w", err)
	}
	if !authorized {
		return nil, fmt.Errorf("failed to unregister object: %w", types.ErrNotAuthorized)
	}

	ownerRecord, err := engine.GetRelationship(ctx, pol, authRelationship)
	if err != nil {
		return nil, fmt.Errorf("failed to unregister object: %w", err)
	}
	if ownerRecord.Archived {
		return &types.UnregisterObjectResponse{
			Found: true,
		}, nil
	}

	_, err = c.removeObjectRelationships(ctx, engine, pol, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to unregister object: removing orphan relationships: %w", err)
	}

	err = c.archiveObject(ctx, engine, pol, ownerRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to unregister object: archiving object: %w", err)
	}

	return &types.UnregisterObjectResponse{
		Found: true,
	}, nil
}

func (c *UnregisterObjectHandler) archiveObject(ctx context.Context, engine *zanzi.Adapter, policy *types.Policy, ownerRecord *types.RelationshipRecord) error {
	ownerRecord.Archived = true
	_, err := engine.SetRelationship(ctx, policy, ownerRecord)
	return err

}

func (c *UnregisterObjectHandler) removeObjectRelationships(ctx context.Context, engine *zanzi.Adapter, pol *types.Policy, cmd *types.UnregisterObjectRequest) (uint, error) {
	selector := &types.RelationshipSelector{
		ObjectSelector: &types.ObjectSelector{
			Selector: &types.ObjectSelector_Object{
				Object: cmd.Object,
			},
		},
		RelationSelector: &types.RelationSelector{
			Selector: &types.RelationSelector_Wildcard{
				Wildcard: &types.WildcardSelector{},
			},
		},
		SubjectSelector: &types.SubjectSelector{
			Selector: &types.SubjectSelector_Wildcard{
				Wildcard: &types.WildcardSelector{},
			},
		},
	}
	count, err := engine.DeleteRelationships(ctx, pol, selector)
	if err != nil {
		return 0, fmt.Errorf("failed to unregister object: %w", err)
	}

	return count, nil
}

func (c *UnregisterObjectHandler) validate(ctx context.Context, cmd *types.UnregisterObjectRequest) error {
	if cmd.PolicyId == "" {
		return types.ErrPolicyNil
	}

	if cmd.Object == nil {
		return types.ErrObjectNil
	}

	return nil
}
