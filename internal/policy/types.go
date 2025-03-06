package policy

import (
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

// PolicyIR is an intermediary representation of a Policy which marshaled representations
// must unmarshall to.
type PolicyIR struct {
	Name          string
	Description   string
	Attributes    map[string]string
	Resources     []*types.Resource
	ActorResource *types.ActorResource
}

// sort performs an in place sorting of resources, relations and permissions in a policy
func (pol *PolicyIR) Sort() {
	resourceExtractor := func(resource *types.Resource) string { return resource.Name }
	relationExtractor := func(relation *types.Relation) string { return relation.Name }
	permissionExtractor := func(permission *types.Permission) string { return permission.Name }

	utils.FromExtractor(pol.Resources, resourceExtractor).SortInPlace()

	for _, resource := range pol.Resources {
		utils.FromExtractor(resource.Relations, relationExtractor).SortInPlace()
		utils.FromExtractor(resource.Permissions, permissionExtractor).SortInPlace()
	}
}
