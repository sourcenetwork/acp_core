package relationship

import (
	"context"

	"github.com/sourcenetwork/acp_core/internal/authorizer"
	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/internal/zanzi"
	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/did"
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
	pol := rec.Policy

	err = c.validate(cmd)
	if err != nil {
		return nil, newRegisterObjectErr(err)
	}

	err = nil

	record, err := queryOwnerRelationship(ctx, engine, pol, registration.Object)
	if err != nil {
		return nil, newRegisterObjectErr(err)
	}
	if record != nil {
		return nil, errors.Wrap("object already registered", errors.ErrorType_OPERATION_FORBIDDEN,
			errors.Pair("policy", cmd.PolicyId),
			errors.Pair("resource", cmd.Object.Resource),
			errors.Pair("id", cmd.Object.Id),
		)
	}

	record = &types.RelationshipRecord{
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
		CreationTime: cmd.CreationTime,
		Metadata:     cmd.Metadata,
	}
	_, err = engine.SetRelationship(ctx, pol, record)
	if err != nil {
		return nil, newRegisterObjectErr(err)
	}

	// TODO efactor the event type
	eventManager.EmitEvent(&types.EventObjectRegistered{
		Actor:          registration.Actor.Id,
		PolicyId:       pol.Id,
		ObjectResource: registration.Object.Resource,
		ObjectId:       registration.Object.Id,
	})

	return &types.RegisterObjectResponse{
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

type ArchiveObjectHandler struct{}

func (c *ArchiveObjectHandler) Execute(ctx context.Context, runtime runtime.RuntimeManager, cmd *types.ArchiveObjectRequest) (*types.ArchiveObjectResponse, error) {
	err := c.validateCmd(cmd)
	if err != nil {
		return nil, newArchiveObjectErr(err)
	}

	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, newArchiveObjectErr(err)
	}

	principal, err := auth.ExtractPrincipalWithType(ctx, auth.DID)
	if err != nil {
		return nil, newArchiveObjectErr(err)
	}
	did := principal.Identifier()

	rec, err := engine.GetPolicy(ctx, cmd.PolicyId)
	if err != nil {
		return nil, newArchiveObjectErr(err)
	}
	if rec == nil {
		return nil, newArchiveObjectErr(errors.NewPolicyNotFound(cmd.PolicyId))
	}
	pol := rec.Policy

	ownerRecord, err := queryOwnerRelationship(ctx, engine, pol, cmd.Object)
	if err != nil {
		return nil, newArchiveObjectErr(err)
	}
	if ownerRecord == nil {
		return nil, newArchiveObjectErr(errors.Wrap("object not registered", errors.ErrorType_BAD_INPUT,
			errors.Pair("policy", pol.Id),
			errors.Pair("resource", cmd.Object.Resource),
			errors.Pair("id", cmd.Object.Id),
		))
	}
	if ownerRecord.Archived {
		return &types.ArchiveObjectResponse{
			RelationshipsRemoved: 0,
			RecordModified:       false,
		}, nil // noop when object is archived
	}

	authorizer := authorizer.NewOperationAuthorizer(engine)

	mutateOwnerOperation := types.Operation{
		Object:     cmd.Object,
		Permission: policy.OwnerRelation,
	}
	actor := types.Actor{Id: did}
	authorized, err := authorizer.IsAuthorized(ctx, pol, &mutateOwnerOperation, &actor)
	if err != nil {
		return nil, newArchiveObjectErr(err)
	}
	if !authorized {
		return nil, newArchiveObjectErr(errors.Wrap("actor cannot manage owner relation", errors.ErrorType_UNAUTHORIZED,
			errors.Pair("policy", cmd.PolicyId),
			errors.Pair("resource", cmd.Object.Resource),
			errors.Pair("object", cmd.Object.Id),
			errors.Pair("actor", did),
		))
	}

	count, err := c.removeObjectRelationships(ctx, engine, pol, cmd)
	if err != nil {
		return nil, newArchiveObjectErr(err)
	}

	ownerRecord.Archived = true
	_, err = engine.SetRelationship(ctx, pol, ownerRecord)
	if err != nil {
		return nil, errors.Wrap("archiving object", err,
			errors.Pair("policy", ownerRecord.PolicyId),
			errors.Pair("resource", ownerRecord.Relationship.Object.Resource),
			errors.Pair("object", ownerRecord.Relationship.Object.Id),
		)
	}

	return &types.ArchiveObjectResponse{
		RelationshipsRemoved: count,
		RecordModified:       true,
	}, nil
}

func (c *ArchiveObjectHandler) removeObjectRelationships(ctx context.Context, engine *zanzi.Adapter, pol *types.Policy, cmd *types.ArchiveObjectRequest) (uint64, error) {
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

func (c *ArchiveObjectHandler) validateCmd(cmd *types.ArchiveObjectRequest) error {
	if err := ObjectSpec(cmd.Object); err != nil {
		return err
	}
	return nil
}

type TransferObjectHandler struct{}

func (h *TransferObjectHandler) Execute(ctx context.Context, runtime runtime.RuntimeManager, cmd *types.TransferObjectRequest) (*types.TransferObjectResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, newTransferObjectErr(err)
	}

	principal, err := auth.ExtractPrincipalWithType(ctx, auth.DID)
	if err != nil {
		return nil, newTransferObjectErr(err)
	}
	did := principal.Identifier()

	rec, err := engine.GetPolicy(ctx, cmd.PolicyId)
	if err != nil {
		return nil, newTransferObjectErr(err)
	}
	if rec == nil {
		return nil, newTransferObjectErr(errors.NewPolicyNotFound(cmd.PolicyId))
	}
	pol := rec.Policy

	ownerRecord, err := queryOwnerRelationship(ctx, engine, pol, cmd.Object)
	if err != nil {
		return nil, newTransferObjectErr(err)
	}
	if ownerRecord == nil {
		return nil, newTransferObjectErr(errors.Wrap("cannot transfer an unregistered object",
			errors.ErrorType_NOT_FOUND,
			errors.Pair("policy", pol.Id),
			errors.Pair("resource", cmd.Object.Resource),
			errors.Pair("object", cmd.Object.Id),
		))
	}

	operation := types.Operation{
		Object:     cmd.Object,
		Permission: policy.OwnerRelation,
	}
	authorizer := authorizer.NewOperationAuthorizer(engine)
	authorized, err := authorizer.IsAuthorized(ctx, pol, &operation, types.NewActor(did))
	if err != nil {
		return nil, newTransferObjectErr(err)
	}
	if !authorized {
		return nil, newTransferObjectErr(errors.Wrap("cannot transfer object owned by someone else",
			errors.ErrorType_UNAUTHORIZED,
			errors.Pair("policy", pol.Id),
			errors.Pair("resource", cmd.Object.Resource),
			errors.Pair("object", cmd.Object.Id),
		))
	}

	_, err = engine.DeleteRelationship(ctx, pol, ownerRecord.Relationship)
	if err != nil {
		return nil, newTransferObjectErr(err)
	}

	ownerRecord.Relationship.Subject = &types.Subject{
		Subject: &types.Subject_Actor{
			Actor: cmd.NewOwner,
		},
	}
	_, err = engine.SetRelationship(ctx, pol, ownerRecord)
	if err != nil {
		return nil, newTransferObjectErr(err)
	}

	return &types.TransferObjectResponse{
		Record: ownerRecord,
	}, nil
}

type AmendRegistrationHandler struct{}

func (h *AmendRegistrationHandler) Handle(ctx context.Context, runtime runtime.RuntimeManager, req *types.AmendRegistrationRequest) (*types.AmendRegistrationResponse, error) {
	_, err := auth.ExtractPrincipalWithType(ctx, auth.Root)
	if err != nil {
		return nil, newAmendRegistrationErr(err)
	}

	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, newAmendRegistrationErr(err)
	}

	policy, relRec, err := h.verifyPreconditions(ctx, engine, req)
	if err != nil {
		return nil, newAmendRegistrationErr(err)
	}

	_, err = engine.DeleteRelationship(ctx, policy, relRec.Relationship)
	if err != nil {
		return nil, newAmendRegistrationErr(errors.Wrap("removing old relationship", err))
	}

	relRec.OwnerDid = req.NewOwner.Id
	relRec.Relationship.Subject = &types.Subject{
		Subject: &types.Subject_Actor{
			Actor: req.NewOwner,
		},
	}

	_, err = engine.SetRelationship(ctx, policy, relRec)
	if err != nil {
		return nil, newAmendRegistrationErr(errors.Wrap("creating new relationship", err))
	}

	return &types.AmendRegistrationResponse{
		Record: relRec,
	}, nil
}

func (h *AmendRegistrationHandler) verifyPreconditions(ctx context.Context, engine *zanzi.Adapter, req *types.AmendRegistrationRequest) (*types.Policy, *types.RelationshipRecord, error) {
	polRec, err := engine.GetPolicy(ctx, req.PolicyId)
	if err != nil {
		return nil, nil, err
	}
	if polRec == nil {
		return nil, nil, errors.NewPolicyNotFound(req.PolicyId)
	}

	err = did.IsValidDID(req.NewOwner.Id)
	if err != nil {
		return nil, nil, errors.NewFromBaseError(err, errors.ErrorType_BAD_INPUT, "invalid actor id", errors.Pair("id", req.NewOwner.Id))
	}

	relRec, err := queryOwnerRelationship(ctx, engine, polRec.Policy, req.Object)
	if err != nil {
		return nil, nil, err
	}
	if relRec == nil {
		return nil, nil, errors.Wrap("object not registered", errors.ErrorType_BAD_INPUT,
			errors.Pair("policy", req.PolicyId),
			errors.Pair("resource", req.Object.Resource),
			errors.Pair("id", req.Object.Id),
		)
	}

	return polRec.Policy, relRec, nil
}

type UnarchiveObjectHandler struct{}

func (h *UnarchiveObjectHandler) Handle(ctx context.Context, runtime runtime.RuntimeManager, cmd *types.UnarchiveObjectRequest) (*types.UnarchiveObjectResponse, error) {
	engine, err := zanzi.NewZanzi(runtime.GetKVStore(), runtime.GetLogger())
	if err != nil {
		return nil, newRegisterObjectErr(err)
	}

	rec, err := engine.GetPolicy(ctx, cmd.PolicyId)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, newRegisterObjectErr(errors.NewPolicyNotFound(cmd.PolicyId))
	}
	pol := rec.Policy

	principal, err := auth.ExtractPrincipalWithType(ctx, auth.DID)
	if err != nil {
		return nil, newRegisterObjectErr(err)
	}
	did := principal.Identifier()

	ownerRecord, err := queryOwnerRelationship(ctx, engine, pol, cmd.Object)
	if err != nil {
		return nil, newArchiveObjectErr(err)
	}
	if ownerRecord == nil {
		return nil, newUnarchiveObjectErr(errors.Wrap("object not registered", errors.ErrorType_BAD_INPUT,
			errors.Pair("policy", pol.Id),
			errors.Pair("resource", cmd.Object.Resource),
			errors.Pair("id", cmd.Object.Id),
		))
	}
	if !ownerRecord.Archived {
		return &types.UnarchiveObjectResponse{
			Record:         ownerRecord,
			RecordModified: false,
		}, nil
	}

	authorizer := authorizer.NewOperationAuthorizer(engine)

	mutateOwnerOperation := types.Operation{
		Object:     cmd.Object,
		Permission: policy.OwnerRelation,
	}
	actor := types.Actor{Id: did}
	authorized, err := authorizer.IsAuthorized(ctx, pol, &mutateOwnerOperation, &actor)
	if err != nil {
		return nil, newUnarchiveObjectErr(err)
	}
	if !authorized {
		return nil, newUnarchiveObjectErr(errors.Wrap("actor cannot manage owner relation", errors.ErrorType_UNAUTHORIZED,
			errors.Pair("policy", cmd.PolicyId),
			errors.Pair("resource", cmd.Object.Resource),
			errors.Pair("object", cmd.Object.Id),
			errors.Pair("actor", did),
		))
	}

	ownerRecord.Archived = false
	_, err = engine.SetRelationship(ctx, pol, ownerRecord)
	if err != nil {
		return nil, errors.Wrap("updating record", err,
			errors.Pair("policy", ownerRecord.PolicyId),
			errors.Pair("resource", ownerRecord.Relationship.Object.Resource),
			errors.Pair("object", ownerRecord.Relationship.Object.Id),
		)
	}

	return &types.UnarchiveObjectResponse{
		Record:         ownerRecord,
		RecordModified: false,
	}, nil
}
