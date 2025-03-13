package ppp

import (
	"github.com/sourcenetwork/acp_core/pkg/transformer"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func PipelineFactory(sequence uint64, spec types.PolicySpecification) Pipeline {
	var transformers []transformer.Transformer
	var specs []transformer.Specification
	switch spec {
	case types.PolicySpecification_DEFRA_SPEC:
		transformers = []transformer.Transformer{}
		specs = []transformer.Specification{
			transformer.NewDefraSpec(),
		}
	case types.PolicySpecification_UNKNOWN_SPEC:
	case types.PolicySpecification_NO_SPEC:
	default:
	}
	return newPipeline(sequence, specs, transformers)
}
