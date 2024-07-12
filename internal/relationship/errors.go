package relationship

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
)

var (
	ErrDeleteOwnerRel = errors.Wrap("deleting an owner relationship is not allowed, consider archiving the object", errors.ErrorType_OPERATION_FORBIDDEN)
	ErrSetOwnerRel    = errors.Wrap("creating an owner relationship is not allowed, consider registering the object", errors.ErrorType_OPERATION_FORBIDDEN)
)

func newFilterRelationshpErr(err error) error {
	return errors.Wrap("filter relationships error", err)
}

func newSetRelationshipErr(err error) error {
	return errors.Wrap("set relationship error", err)
}

func newDeleteRelationshipErr(err error) error {
	return errors.Wrap("delete relationship erro", err)
}

func newGetObjectRegistrationErr(err error) error {
	return errors.Wrap("get object error", err)
}

func newRegisterObjectErr(err error) error {
	return errors.Wrap("register object error", err)
}

func newUnregisterObjectErr(err error) error {
	return errors.Wrap("unregister object error", err)
}
