package theorem

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
)

func newEvaluatorErr(err error) error {
	return errors.Wrap("theorem evaluator failed", err)
}
