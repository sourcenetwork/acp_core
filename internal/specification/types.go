// package specification provides hook abstractions which can be used to transform and validate policies
package specification

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// Requirement models a component which can validate that a Policy meets some arbitrary specification.
//
// The Requirement can be caller implemented, which effectively constrains the Policies that acp_core
// understands as valid.
type Requirement interface {
	// Validate executes arbitary validation code against the given Policy
	// Returns all returns found within Policy as a MultiError
	Validate(policy types.Policy) *errors.MultiError

	// GetBaseError returns the base error used to construct Specification errors
	GetBaseError() error
}

// Transformer models a component which transforms the given policy into another one
//
// This abstraction allows callers to hook into the policy processing pipeline system
// within ACP core and add custom validation to it
//
// Transformer extends the Specification interface in order to add redundancy validation to a transformation,
// since it's possible that a subsequent transformation undoes the transformation done by a previous one
type Transformer interface {
	Requirement

	// Transforms takes as input a Policy and maps into another Policy
	Transform(policy types.Policy) (types.Policy, error)
}

// Specification models a set of criteria a Policy must match
// to satisfy some specification type
type Specification interface {
	// GetType returns the known type for the current specification
	// Returns UNKOWN specification if not registered
	GetType() types.PolicySpecificationType

	// GetRequirements return the set of requirements a Policy
	// must satisfy in order to satisfy the specification
	GetRequirements() []Requirement

	// GetTransformers returns the set of transform operations
	// which must be applied to a policy of a given specification type
	GetTransformers() []Transformer
}

var _ Specification = (*specification)(nil)

// specification implements the Specification interface
type specification struct {
	specType     types.PolicySpecificationType
	requirements []Requirement
	transformers []Transformer
}

func (s *specification) GetType() types.PolicySpecificationType {
	return s.specType
}

func (s *specification) GetTransformers() []Transformer {
	return s.transformers
}

func (s *specification) GetRequirements() []Requirement {
	return s.requirements
}

// newSpecification returns an instance of Specification with the given type, requirements and transformers
func newSpecification(t types.PolicySpecificationType, requirements []Requirement, transformers []Transformer) Specification {
	return &specification{
		specType:     t,
		requirements: requirements,
		transformers: transformers,
	}
}

func NoSpecification() Specification {
	return newSpecification(types.PolicySpecificationType_NO_SPEC, nil, nil)
}
