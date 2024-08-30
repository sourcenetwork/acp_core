//go:build js

package playground_js

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
)

func newInvalidArgsErr(count int) error {
	return errors.Wrap("invalid number of arguments: expected 1", errors.ErrorType_BAD_INPUT, errors.Pair("args", count))
}
