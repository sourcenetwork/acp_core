package ppp

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	DefraReadPermission  = "read"
	DefraWritePermission = "write"
)

var RequiredPermissions = sets.New(DefraReadPermission, DefraWritePermission)

var _ Specification = (*DefraSpec)(nil)

type DefraSpec struct{}

func (s *DefraSpec) Validate(pol *types.Policy) *errors.MultiError {
	// for every resource, there exists permissions read and write
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

func (s *DefraSpec) Name() string {
	return "Defra Specification"
}
