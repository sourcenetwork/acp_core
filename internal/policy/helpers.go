package policy

import (
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func BuildManagementGraph(policy *types.Policy) *types.ManagementGraph {
	graph := &types.ManagementGraph{}
	graph.LoadFromPolicy(policy)
	return graph
}
