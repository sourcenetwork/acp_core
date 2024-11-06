package relationship

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
)

var (
	ErrDeleteOwnerRel = errors.Wrap("deleting an owner relationship is not allowed, consider archiving the object", errors.ErrorType_OPERATION_FORBIDDEN)
	ErrSetOwnerRel    = errors.Wrap("creating an owner relationship is not allowed, consider registering the object", errors.ErrorType_OPERATION_FORBIDDEN)
)

func newFilterRelationshpErr(err error) error {
	return errors.Wrap("filter relationships failed", err)
}

func newSetRelationshipErr(err error) error {
	return errors.Wrap("set relationship failed", err)
}

func newDeleteRelationshipErr(err error) error {
	return errors.Wrap("delete relationship failed", err)
}

func newGetObjectRegistrationErr(err error) error {
	return errors.Wrap("get object failed", err)
}

func newRegisterObjectErr(err error) error {
	return errors.Wrap("register object failed", err)
}

func newArchiveObjectErr(err error) error {
	return errors.Wrap("archive object failed", err)
}

func newValidateRelationshipErr(err error) error {
	return errors.Wrap("validate relationship failed", err)
}

func newTransferObjectErr(err error) error {
	return errors.Wrap("transfer object failed", err)
}

func newAmendRegistrationErr(err error) error {
	return errors.Wrap("amend registration failed", err)
}

func newUnarchiveObjectErr(err error) error {
	return errors.Wrap("unarchiving object failed", err)
}
