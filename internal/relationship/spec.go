package relationship

import (
	"github.com/sourcenetwork/acp_core/pkg/did"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// relationshipSpec validates a Relationshp according to the expected format.
// For now, we can rely on Zanzi to perform the bulk of the validation,
// however its paramount that a Relationship whose subject type is an Actor must be a DID
func relationshipSpec(pol *types.Policy, relationship *types.Relationship) error {
	switch subj := relationship.Subject.Subject.(type) {
	case *types.Subject_Actor:
		if err := did.IsValidDID(subj.Actor.Id); err != nil {
			return errors.Wrap("actor must be a valid did: "+err.Error(),
				errors.ErrInvalidRelationship, errors.Pair("did", subj.Actor.Id))
		}
	case *types.Subject_Object:
		err := did.IsValidDID(subj.Object.Id)
		if subj.Object.Resource == pol.ActorResource.Name && err != nil {
			return errors.Wrap("actor must be a valid did: "+err.Error(),
				errors.ErrInvalidRelationship, errors.Pair("did", subj.Object.Id))
		}
	}
	if relationship.Object.Id == "" {
		return errors.Wrap("object id must not be empty", errors.ErrInvalidRelationship)
	}
	return nil
}

func registrationSpec(registration *types.Registration) error {
	if registration == nil {
		return errors.Wrap("registration is required", errors.ErrorType_BAD_INPUT)
	}

	if registration.Actor == nil {
		return errors.Wrap("registration actor is required", errors.ErrorType_BAD_INPUT)
	}

	if registration.Object == nil {
		return errors.Wrap("registration object is required", errors.ErrorType_BAD_INPUT)
	}

	if registration.Object.Id == "" {
		return errors.Wrap("registration object id is required", errors.ErrorType_BAD_INPUT)
	}

	if err := did.IsValidDID(registration.Actor.Id); err != nil {
		return errors.Wrap("invalid registration: invalid actor did", errors.ErrorType_BAD_INPUT)
	}

	return nil
}

func ObjectSpec(obj *types.Object) error {
	if obj.Id == "" {
		return errors.Wrap("object ID must not be empty", errors.ErrorType_BAD_INPUT)
	}
	return nil
}
