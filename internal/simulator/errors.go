package simulator

import "github.com/sourcenetwork/acp_core/pkg/errors"

func newSimulateError(err error) error {
	return errors.Wrap("Simulate failed", err)
}
