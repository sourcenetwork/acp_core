package ppp

import (
	"github.com/sourcenetwork/acp_core/internal/specification"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func PipelineFactory(sequence uint64, specType types.PolicySpecificationType) Pipeline {
	switch specType {
	case types.PolicySpecificationType_DEFRA_SPEC:
		return newCreatePolicyPipeline(sequence, specification.NewDefraSpecification())
	case types.PolicySpecificationType_UNKNOWN_SPEC:
		return newCreatePolicyPipeline(sequence, specification.NoSpecification())
	case types.PolicySpecificationType_NO_SPEC:
		return newCreatePolicyPipeline(sequence, specification.NoSpecification())
	default:
		return newCreatePolicyPipeline(sequence, specification.NoSpecification())
	}
}
