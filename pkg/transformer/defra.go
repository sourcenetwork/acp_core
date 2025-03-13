package transformer

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
	"k8s.io/apimachinery/pkg/util/sets"
)

var _ Specification = (*DefraSpec)(nil)
var _ Transformer = (*DefraWriteImpliesReadTransform)(nil)

const (
	DefraReadPermissionName  = "read"
	DefraWritePermissionName = "write"
)

// ErrDefraSpec is the base error for the DefraSpec implementation
var ErrDefraSpec = errors.New("defra policy specification", errors.ErrorType_BAD_INPUT)

// RequiredPermissiosn are the set of permissions which all resources
// must include in a Defra compliant Policy.
var RequiredPermissions = sets.New(DefraReadPermissionName, DefraWritePermissionName)

// NewDefraSpec returns an instance of DefraSpec
func NewDefraSpec() Specification {
	return &DefraSpec{}
}

// DefraSpec implements the Specification interface for Defra compliant Policies
//
// Defra compliant Policies require that all resources contain a set of pre-determined
// permissions.
type DefraSpec struct{}

func (s *DefraSpec) Validate(pol types.Policy) *errors.MultiError {
	multiErr := errors.NewMultiError(ErrDefraSpec)

	for _, resource := range pol.Resources {
		permissions := utils.MapSlice(resource.Permissions, func(p *types.Permission) string { return p.Name })
		permissionsSet := sets.New(permissions...)
		intersection := permissionsSet.Intersection(RequiredPermissions)
		if intersection.Len() != RequiredPermissions.Len() {
			missing := RequiredPermissions.Difference(intersection)
			var missingNames []string
			for name, _ := range missing {
				missingNames = append(missingNames, name)
			}
			violation := fmt.Errorf("resource %v: missing permissions: %v", resource.Name, missingNames)
			multiErr.Append(violation)
		}
	}
	if len(multiErr.GetErrors()) > 0 {
		return multiErr
	}
	return nil
}

func (s *DefraSpec) GetBaseError() error { return ErrDefraSpec }

type DefraWriteImpliesReadTransform struct{}

func (s *DefraWriteImpliesReadTransform) GetBaseError() error {
	return ErrDefraSpec
}

func (s *DefraWriteImpliesReadTransform) Validate(pol types.Policy) *errors.MultiError {
	return nil
}

func (s *DefraWriteImpliesReadTransform) Transform(pol types.Policy) (types.Policy, error) {
	return types.Policy{}, nil
}
