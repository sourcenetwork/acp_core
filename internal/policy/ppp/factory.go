package ppp

import (
	"github.com/sourcenetwork/acp_core/internal/specification"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func PipelineFactory(sequence uint64, specType types.PolicySpecificationType) Pipeline {
	switch specType {
	case types.PolicySpecificationType_DEFRA_SPEC:
		return newPipeline(sequence, specification.NewDefraSpecification())
	case types.PolicySpecificationType_UNKNOWN_SPEC:
		return newPipeline(sequence, specification.NoSpecification())
	case types.PolicySpecificationType_NO_SPEC:
		return newPipeline(sequence, specification.NoSpecification())
	default:
		return newPipeline(sequence, specification.NoSpecification())
	}
}
