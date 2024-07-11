package system

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
)

func newSetParamsErr(err error) error {
	return errors.Wrap("set params error", err)
}

func newGetParamsErr(err error) error {
	return errors.Wrap("get params error", err)
}
