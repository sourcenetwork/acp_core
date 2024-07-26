package simulator

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
)

func newSimulateErr(err error) error {
	return errors.Wrap("simulate failed", err)
}
