package policy

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
)

var (
	ErrUnknownMarshalingType = errors.Wrap("unknown marshaling type", errors.ErrorType_BAD_INPUT)
	ErrUnmarshaling          = errors.Wrap("unmarshaling error", errors.ErrorType_BAD_INPUT)

	ErrInvalidShortPolicy           = errors.Wrap("invalid short policy", errors.ErrorType_BAD_INPUT)
	ErrResourceMissingOwnerRelation = errors.Wrap("resource missing owner relation", errors.ErrInvalidPolicy)
	ErrInvalidManagementRule        = errors.Wrap("invalid relation managament definition: %w", errors.ErrInvalidPolicy)
)

func newEvaluateTheoremErr(err error) error {
	return errors.Wrap("evaluate theorem failed", err)
}

func newPolicyCatalogueErr(err error) error {
	return errors.Wrap("get policy catalogue failed", err)
}

func newErrMissingOwnerRelation(resourceName string) error {
	return errors.Wrap("", ErrResourceMissingOwnerRelation,
		errors.Pair("resource", resourceName))
}
