package ppp

import (
	"fmt"

	"github.com/cosmos/gogoproto/proto"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

type Pipeline struct {
	specs        []Specification
	transformers []Transformer
}

// Process executes the Policy Processing Pipeline by sequentially
// applying all transformers (in the order given).
//
// After the transforming step is finished, we iteratte over all transformers and specifcations once again applying the Validate method.
//
// Return the processed policy and an error with the underlying type *errors.MultiError, in case further inspection is necessary.
func (p *Pipeline) Process(pol *types.Policy) (*types.Policy, error) {
	pol, err := p.applyTransforms(pol)
	if err != nil {
		return nil, err
	}

	err = p.applySpecs(pol)
	if err != nil {
		return nil, err
	}

	return pol, nil
}

func (p *Pipeline) applySpecs(pol *types.Policy) *errors.MultiError {
	err := errors.NewMultiError("policy spec verfication", errors.ErrorType_BAD_INPUT)
	specs := make([]Specification, 0, len(p.specs)+len(p.transformers))
	specs = append(specs, p.specs...)
	specs = append(specs, p.transformers...)

	for _, spec := range specs {
		pol := proto.Clone(pol).(*types.Policy)
		violations := spec.Validate(pol)
		utils.MapSlice(violations, func(err error) error {
			return fmt.Errorf("%v: %w", spec.Name(), err)
		})
		err.Append(violations...)
	}
	return err
}

func (p *Pipeline) applyTransforms(policy *types.Policy) (*types.Policy, *errors.MultiError) {
	var err error
	for _, transform := range p.transformers {
		producer := func() *types.Policy {
			return proto.Clone(policy).(*types.Policy)
		}
		policy, err = transform.Transform(producer)
		if err != nil {
			msg := "transform failed: " + transform.Name()
			return nil, errors.NewMultiError(msg, errors.ErrorType_BAD_INPUT, err)
		}
	}
	return policy, nil
}
