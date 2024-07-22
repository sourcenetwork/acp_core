package policy

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

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
func (pol *PolicyIR) sort() {
	resourceExtractor := func(resource *types.Resource) string { return resource.Name }
	relationExtractor := func(relation *types.Relation) string { return relation.Name }
	permissionExtractor := func(permission *types.Permission) string { return permission.Name }

	utils.FromExtractor(pol.Resources, resourceExtractor).SortInPlace()

	for _, resource := range pol.Resources {
		utils.FromExtractor(resource.Relations, relationExtractor).SortInPlace()
		utils.FromExtractor(resource.Permissions, permissionExtractor).SortInPlace()
	}
}

// policyIder builds Policy ids
type policyIder struct{}

// buildId computes the unique id for a policy.
//
// the id is a hash of the policy hash and the policy counter number
func (i *policyIder) Id(pol *types.Policy, counter uint64) string {
	hasher := sha256.New()

	hasher.Write(i.hashPol(pol))
	hasher.Write([]byte(fmt.Sprintf("%v", counter)))

	hash := hasher.Sum(nil)
	id := hex.EncodeToString(hash)
	return id
}

// hashPol computes a partial sha256 hash of a policy.
// the hashing algorithm includes a subset of the fields which are deterministic.
func (i *policyIder) hashPol(pol *types.Policy) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(pol.Name))

	for _, resource := range pol.Resources {
		hasher.Write([]byte(resource.Name))

		for _, rel := range resource.Relations {
			hasher.Write([]byte(rel.Name))
		}

		for _, perm := range resource.Permissions {
			hasher.Write([]byte(perm.Name))
			hasher.Write([]byte(perm.Expression))
		}
	}

	return hasher.Sum(nil)
}
