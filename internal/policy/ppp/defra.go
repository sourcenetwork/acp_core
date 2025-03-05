package ppp

import (
	"fmt"

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

func (s *DefraSpec) Validate(pol *types.Policy) []error {
	// for every resource, there exists permissions read and write
	var violations []error
	for _, resource := range pol.Resources {
		permissions := utils.MapSlice(resource.Permissions, func(p *types.Permission) string { return p.Name })
		permissionsSet := sets.New(permissions...)
		intersection := permissionsSet.Intersection(RequiredPermissions)
		if intersection.Len() != RequiredPermissions.Len() {
			missing := RequiredPermissions.Difference(intersection)
			violation := fmt.Errorf("resource %v: missing permissions: %v", resource.Name, missing)
			violations = append(violations, violation)
		}
	}
	return violations
}

func (s *DefraSpec) Name() string {
	return "Defra Specification"
}
