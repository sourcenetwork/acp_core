package ppp

import (
	"github.com/cosmos/gogoproto/proto"
	"github.com/sourcenetwork/acp_core/internal/specification"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var ErrPolicyProcessing = errors.New("policy processing", errors.ErrorType_BAD_INPUT)

var createPolicyRequirements = []specification.Requirement{
	&BasicRequirement{},
}

func newCreatePolicyPipeline(sequenceNumber uint64, spec specification.Specification) Pipeline {
	headTransformers := []specification.Transformer{
		&BasicTransformer{},
		&DiscretionaryTransformer{},
		&DecentralizedAdminTransformer{},
	}

	tailTransformers := []specification.Transformer{
		&SortTransformer{},
		NewIdTransformer(sequenceNumber),
	}

	transformerPipeline := headTransformers
	transformerPipeline = append(transformerPipeline, spec.GetTransformers()...)
	transformerPipeline = append(transformerPipeline, tailTransformers...)

	requirements := append(createPolicyRequirements, spec.GetRequirements()...)

	return Pipeline{
		requirements: requirements,
		transformers: transformerPipeline,
	}
}

// newEditPolicyPipeline returns a Pipeline which handles transforms while editing a Policy
func newEditPolicyPipeline(oldPolicy *types.Policy, spec specification.Specification) Pipeline {
	requirements := []specification.Requirement{
		&BasicRequirement{},
		NewImmutableSpecRequirement(oldPolicy.SpecificationType),
		NewPreservedResourcesRequirement(oldPolicy),
	}

	headTransformers := []specification.Transformer{
		&BasicTransformer{},
		NewTransferIdTransformer(oldPolicy.Id),
		&DiscretionaryTransformer{},
		&DecentralizedAdminTransformer{},
	}

	tailTransformers := []specification.Transformer{
		&SortTransformer{},
	}

	transformerPipeline := headTransformers
	transformerPipeline = append(transformerPipeline, spec.GetTransformers()...)
	transformerPipeline = append(transformerPipeline, tailTransformers...)

	requirements = append(requirements, spec.GetRequirements()...)

	return Pipeline{
		requirements: requirements,
		transformers: transformerPipeline,
	}
}

type Pipeline struct {
	requirements []specification.Requirement
	transformers []specification.Transformer
}

// Process executes the Policy Processing Pipeline by sequentially
// applying all transformers (in the order given).
//
// After the transforming step is finished, we iteratte over all transformers and specifcations once again applying the Validate method.
//
// Return the processed policy and an error with the underlying type *errors.MultiError, in case further inspection is necessary.
func (p *Pipeline) Process(pol *types.Policy) (specification.TransformerResult, error) {
	result, err := p.applyTransforms(*pol)
	if err != nil {
		return specification.TransformerResult{}, err
	}

	multiErr := p.validateRequirements(&result.Policy)
	if multiErr != nil {
		return specification.TransformerResult{}, multiErr
	}

	return result, nil
}

func (p *Pipeline) validateRequirements(pol *types.Policy) *errors.MultiError {
	multiErr := errors.NewMultiError(ErrPolicyProcessing)

	for _, spec := range p.requirements {
		// clone each policy before sending to the spec to ensure it's
		// a buggy Requirement doesn't ruin the Policy
		clone := proto.Clone(pol).(*types.Policy)
		err := spec.Validate(*clone)
		if err != nil {
			multiErr.Append(err)
		}
	}

	for _, transformer := range p.transformers {
		cloned := proto.Clone(pol).(*types.Policy)
		err := transformer.Validate(*cloned)
		if err != nil {
			multiErr.Append(err)
		}
	}

	if len(multiErr.GetErrors()) == 0 {
		return nil
	}

	return multiErr
}

func (p *Pipeline) applyTransforms(policy types.Policy) (specification.TransformerResult, error) {
	result := specification.TransformerResult{}
	result.Policy = policy
	for _, trans := range p.transformers {
		tmp, err := trans.Transform(result.Policy)
		result.Policy = tmp.Policy
		result.Messages = append(result.Messages, tmp.Messages...)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}
