package policy

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
)

var (
	ErrUnknownMarshalingType = errors.Wrap("unknown marshaling type", errors.ErrorType_BAD_INPUT)
	ErrUnmarshaling          = errors.Wrap("unmarshaling error", errors.ErrorType_BAD_INPUT)

	ErrInvalidShortPolicy = errors.Wrap("invalid short policy", errors.ErrorType_BAD_INPUT)
	ErrInvalidYamlPolicy  = errors.Wrap("invalid yaml policy", errors.ErrorType_BAD_INPUT)

	ErrResourceNotInPolicy     = errors.Wrap("resource not in policy", errors.ErrorType_BAD_INPUT)
	ErrRelationNotInResource   = errors.Wrap("relation not in resource", errors.ErrorType_BAD_INPUT)
	ErrPermissionNotInResource = errors.Wrap("permission not in resource", errors.ErrorType_BAD_INPUT)
)

func newEvaluateTheoremErr(err error) error {
	return errors.Wrap("evaluate theorem failed", err)
}

func newPolicyCatalogueErr(err error) error {
	return errors.Wrap("get policy catalogue failed", err)
}

func NewErrResourceNotInPolicy(policyId string, resource string) error {
	return errors.Wrap("", ErrResourceNotInPolicy, errors.Pair("policy", policyId), errors.Pair("resource", resource))
}

func NewErrRelationNotInResource(policyId, resource, relation string) error {
	return errors.Wrap("", ErrRelationNotInResource, errors.Pair("policy", policyId), errors.Pair("resource", resource), errors.Pair("relation", relation))
}

func NewErrPermissionNotInResource(policyId, resource, permission string) error {
	return errors.Wrap("", ErrRelationNotInResource, errors.Pair("policy", policyId), errors.Pair("resource", resource), errors.Pair("permission", permission))
}
