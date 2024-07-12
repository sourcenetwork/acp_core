package relationship

import (
	"context"

	prototypes "github.com/cosmos/gogoproto/types"

	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/errors"
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
		return nil, newRegisterObjectErr(err)
	}

	principal, err := auth.ExtractPrincipalWithType(ctx, auth.DID)
	if err != nil {
		return nil, newRegisterObjectErr(err)
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
		return nil, newRegisterObjectErr(errors.NewPolicyNotFound(cmd.PolicyId))
	}
	policy := rec.Policy

	err = c.validate(cmd)
	if err != nil {
		return nil, newRegisterObjectErr(err)
	}

	err = nil
	var result types.RegistrationResult

	record, err := queryOwnerRelationship(ctx, engine, policy, registration.Object)
	if err != nil {
		return nil, newRegisterObjectErr(err)
	}

	switch c.resolveObjectStatus(record) {
	case statusUnregistered:
		result, err = c.unregisteredStrategy(ctx, engine, policy, registration, did, cmd.CreationTime, cmd.Metadata)
	case statusArchived:
		result, err = c.archivedObjectStrategy(ctx, engine, policy, record, registration)
	case statusActive:
		result, err = c.activeObjectStrategy(record, registration)
	}

	if err != nil {
		return nil, newRegisterObjectErr(err)
	}

	record, err = queryOwnerRelationship(ctx, engine, policy, registration.Object)
	if err != nil {
		return nil, newRegisterObjectErr(err)
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
func (c *RegisterObjectHandler) validate(cmd *types.RegisterObjectRequest) error {
	if err := ObjectSpec(cmd.Object); err != nil {
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

// unregisteredStrategy creates a relationship with the relation `owner` for the object in Registration
func (c *RegisterObjectHandler) unregisteredStrategy(ctx context.Context, zanzi *zanzi.Adapter, pol *types.Policy, registration *types.Registration, creator string, creationTs *prototypes.Timestamp, metadata map[string]string) (types.RegistrationResult, error) {
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
		OwnerDid:     registration.Actor.Id,
		PolicyId:     pol.Id,
		Archived:     false,
		CreationTime: creationTs,
		Metadata:     metadata,
	}
	_, err := zanzi.SetRelationship(ctx, pol, &record)
	if err != nil {
		return types.RegistrationResult_NoOp, err
	}

	return types.RegistrationResult_Registered, nil
}

func (c *RegisterObjectHandler) activeObjectStrategy(record *types.RelationshipRecord, registration *types.Registration) (types.RegistrationResult, error) {
	if record.OwnerDid != registration.Actor.Id {
		return types.RegistrationResult_NoOp, errors.Wrap("object is already registered to a different actor", errors.ErrorType_UNAUTHORIZED,
			errors.Pair("policy", record.PolicyId),
			errors.Pair("resource", record.Relationship.Object.Resource),
			errors.Pair("object", record.Relationship.Object.Id),
			errors.Pair("owner", record.OwnerDid),
		)
	}
	return types.RegistrationResult_NoOp, nil
}

// archivedObjectStrategy modifies the relationship record to be unarchived
func (c *RegisterObjectHandler) archivedObjectStrategy(ctx context.Context, engine *zanzi.Adapter, policy *types.Policy, record *types.RelationshipRecord, registration *types.Registration) (types.RegistrationResult, error) {
	if record.OwnerDid != registration.Actor.Id {
		return types.RegistrationResult_NoOp, errors.Wrap("object was archived by a different actor", errors.ErrorType_UNAUTHORIZED,
			errors.Pair("policy", record.PolicyId),
			errors.Pair("resource", record.Relationship.Object.Resource),
			errors.Pair("object", record.Relationship.Object.Id),
			errors.Pair("owner", record.OwnerDid),
		)
	}

	record.Archived = false
	_, err := engine.SetRelationship(ctx, policy, record)
	if err != nil {
		return types.RegistrationResult_NoOp, err
	}

	return types.RegistrationResult_Unarchived, nil
}

type UnregisterObjectHandler struct{}

func (c *UnregisterObjectHandler) Execute(ctx context.Context, runtime runtime.RuntimeManager, cmd *types.UnregisterObjectRequest) (*types.UnregisterObjectResponse, error) {
	err := c.validateCmd(cmd)
	if err != nil {
		return nil, newUnregisterObjectErr(err)
	}

	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, newUnregisterObjectErr(err)
	}

	principal, err := auth.ExtractPrincipalWithType(ctx, auth.DID)
	if err != nil {
		return nil, newUnregisterObjectErr(err)
	}
	did := principal.Identifier()

	rec, err := engine.GetPolicy(ctx, cmd.PolicyId)
	if err != nil {
		return nil, newUnregisterObjectErr(err)
	}
	if rec == nil {
		return nil, newUnregisterObjectErr(errors.NewPolicyNotFound(cmd.PolicyId))
	}
	pol := rec.Policy

	ownerRecord, err := queryOwnerRelationship(ctx, engine, pol, cmd.Object)
	if err != nil {
		return nil, newUnregisterObjectErr(err)
	}
	if ownerRecord == nil {
		return &types.UnregisterObjectResponse{
			Found:                false,
			RelationshipsRemoved: 0,
		}, nil //noop when object does not exist
	}
	if ownerRecord.Archived {
		return &types.UnregisterObjectResponse{
			Found:                true,
			RelationshipsRemoved: 0,
		}, nil // noop when object is archived
	}

	authorizer := NewRelationshipAuthorizer(engine)

	authRelationship := types.NewActorRelationship(cmd.Object.Resource, cmd.Object.Id, policy.OwnerRelation, did)
	actor := types.Actor{Id: did}
	authorized, err := authorizer.IsAuthorized(ctx, pol, authRelationship, &actor)
	if err != nil {
		return nil, newUnregisterObjectErr(err)
	}
	if !authorized {
		return nil, newUnregisterObjectErr(errors.Wrap("cannot unregister object: actor is not the owner", errors.ErrorType_UNAUTHORIZED,
			errors.Pair("policy", cmd.PolicyId),
			errors.Pair("resource", cmd.Object.Resource),
			errors.Pair("object", cmd.Object.Id)))
	}

	count, err := c.removeObjectRelationships(ctx, engine, pol, cmd)
	if err != nil {
		return nil, newUnregisterObjectErr(err)
	}

	err = c.archiveObject(ctx, engine, pol, ownerRecord)
	if err != nil {
		return nil, newUnregisterObjectErr(err)
	}

	return &types.UnregisterObjectResponse{
		Found:                true,
		RelationshipsRemoved: count,
	}, nil
}

func (c *UnregisterObjectHandler) archiveObject(ctx context.Context, engine *zanzi.Adapter, policy *types.Policy, ownerRecord *types.RelationshipRecord) error {
	ownerRecord.Archived = true
	_, err := engine.SetRelationship(ctx, policy, ownerRecord)
	if err != nil {
		return errors.Wrap("archiving object", err,
			errors.Pair("policy", ownerRecord.PolicyId),
			errors.Pair("resource", ownerRecord.Relationship.Object.Resource),
			errors.Pair("object", ownerRecord.Relationship.Object.Id),
		)
	}
	return nil
}

func (c *UnregisterObjectHandler) removeObjectRelationships(ctx context.Context, engine *zanzi.Adapter, pol *types.Policy, cmd *types.UnregisterObjectRequest) (uint64, error) {
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
		return 0, errors.Wrap("could not delete associated relationships", err,
			errors.Pair("policy", pol.Id),
			errors.Pair("resource", cmd.Object.Resource),
			errors.Pair("object", cmd.Object.Id),
		)
	}
	return uint64(count), nil
}

func (c *UnregisterObjectHandler) validateCmd(cmd *types.UnregisterObjectRequest) error {
	if err := ObjectSpec(cmd.Object); err != nil {
		return err
	}
	return nil
}
