package ppp

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/specification"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
	"k8s.io/apimachinery/pkg/util/sets"
)

var _ specification.Requirement = (*ImmutableIdRequirement)(nil)
var _ specification.Requirement = (*ImmutableSpecRequirement)(nil)
var _ specification.Requirement = (*PreservedResourcesRequirement)(nil)

var ErrPreserveResource = errors.New("cannot remove resources from policy", errors.ErrorType_BAD_INPUT)
var ErrImmutableId = errors.New("editing policy must preserve id", errors.ErrorType_BAD_INPUT)
var ErrImmutableSpec = errors.New("editing policy must preserve spec", errors.ErrorType_BAD_INPUT)

// NewImmutableIdRequirement returns an instance of ImmutableIdRequirement
// bound to some id
func NewImmutableIdRequirement(oldId string) specification.Requirement {
	return &ImmutableIdRequirement{
		oldId: oldId,
	}
}

// ImmutableIdRequirement validates that the new policy has the same Id as the old one
type ImmutableIdRequirement struct {
	oldId string
}

func (r *ImmutableIdRequirement) Validate(policy types.Policy) *errors.MultiError {
	if policy.Id != r.oldId {
		return errors.NewMultiError(ErrImmutableId)
	}
	return nil
}

func (r *ImmutableIdRequirement) GetBaseError() error {
	return ErrImmutableId
}

// NewImmutableSpecRequirement returns a new instance of ImmutableSpecRequirement
// bound to the given spec
func NewImmutableSpecRequirement(spec types.PolicySpecificationType) specification.Requirement {
	return &ImmutableSpecRequirement{
		oldSpec: spec,
	}
}

// ImmutableSpecRequirement validates that the new policy
// has the same spec as the old one
type ImmutableSpecRequirement struct {
	oldSpec types.PolicySpecificationType
}

func (r *ImmutableSpecRequirement) Validate(policy types.Policy) *errors.MultiError {
	if policy.SpecificationType != r.oldSpec {
		return errors.NewMultiError(ErrImmutableSpec)
	}
	return nil
}

func (r *ImmutableSpecRequirement) GetBaseError() error {
	return ErrImmutableSpec
}

// NewPreservedResourcesRequirement returns an instance of PreservedResourcesRequirement
// which asserts the resources in oldPolicy are preserved
func NewPreservedResourcesRequirement(oldPolicy *types.Policy) specification.Requirement {
	return &PreservedResourcesRequirement{
		oldPolicy: oldPolicy,
	}
}

// PreservedResourcesRequirement validates that the new policy
// does not remove any resource from the old policy
type PreservedResourcesRequirement struct {
	oldPolicy *types.Policy
}

func (r *PreservedResourcesRequirement) Validate(policy types.Policy) *errors.MultiError {
	oldResources := sets.New(r.oldPolicy.ListResourcesNames()...)
	newResources := sets.New(policy.ListResourcesNames()...)
	missing := oldResources.Difference(newResources)
	if missing.Len() == 0 {
		return nil
	}

	errs := utils.MapSlice(missing.UnsortedList(), func(name string) error {
		return fmt.Errorf("removed resource %v", name)
	})

	return errors.NewMultiError(nil, errs...)
}

func (r *PreservedResourcesRequirement) GetBaseError() error {
	return ErrPreserveResource
}
