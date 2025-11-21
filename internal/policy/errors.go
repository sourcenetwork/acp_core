package policy

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
)

var (
	ErrUnknownMarshalingType = errors.Wrap("unknown marshaling type", errors.ErrorType_BAD_INPUT)
	ErrUnmarshaling          = errors.Wrap("unmarshaling error", errors.ErrorType_BAD_INPUT)

	ErrInvalidShortPolicy = errors.Wrap("invalid short policy", errors.ErrorType_BAD_INPUT)
	ErrInvalidYamlPolicy  = errors.Wrap("invalid yaml policy", errors.ErrorType_BAD_INPUT)
)

func newEvaluateTheoremErr(err error) error {
	return errors.Wrap("evaluate theorem failed", err)
}

func newPolicyCatalogueErr(err error) error {
	return errors.Wrap("get policy catalogue failed", err)
}
