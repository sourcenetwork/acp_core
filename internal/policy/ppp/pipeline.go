package ppp

import (
	"github.com/cosmos/gogoproto/proto"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func NewPipeline(specs []Specification, transformers []Transformer) Pipeline {
	return Pipeline{
		specs:        specs,
		transformers: transformers,
	}
}

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
	multiErr := errors.NewMultiError(ErrPolicyProcessing)

	for _, spec := range p.specs {
		pol := proto.Clone(pol).(*types.Policy)
		err := spec.Validate(pol)
		if err != nil {
			multiErr.Append(err)
		}
	}

	for _, transformer := range p.transformers {
		pol := proto.Clone(pol).(*types.Policy)
		err := transformer.Validate(pol)
		if err != nil {
			multiErr.Append(err)
		}
	}
	if len(multiErr.GetErrors()) == 0 {
		return nil
	}

	return multiErr
}

func (p *Pipeline) applyTransforms(policy *types.Policy) (*types.Policy, *errors.MultiError) {
	var err *errors.MultiError
	for _, transform := range p.transformers {
		producer := func() *types.Policy {
			return proto.Clone(policy).(*types.Policy)
		}
		policy, err = transform.Transform(producer)
		if err != nil {
			return nil, err
		}
	}
	return policy, nil
}
