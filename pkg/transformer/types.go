// package transformer provides hook abstractions which can be used to transform and validate policies
package transformer

import (
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// Specification models a component which can validate that a Policy meets some arbitrary specification.
//
// The Specification can be caller implemented, which effectively constrains the Policies that acp_core
// understands as valid.
type Specification interface {
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
	Specification

	// Transforms takes as input a Policy and maps into another Policy
	Transform(policy types.Policy) (types.Policy, error)
}
