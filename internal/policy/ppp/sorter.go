package ppp

import (
	"github.com/sourcenetwork/acp_core/internal/specification"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

var _ specification.Transformer = (*SortTransformer)(nil)

// SortTransformer performs a stable sorts over all resources, permissions and relations by their name.
type SortTransformer struct{}

func (t *SortTransformer) Validate(_ types.Policy) *errors.MultiError {
	return nil
}

func (t *SortTransformer) Transform(pol types.Policy) (specification.TransformerResult, error) {
	resourceExtractor := func(resource *types.Resource) string { return resource.Name }
	relationExtractor := func(relation *types.Relation) string { return relation.Name }
	permissionExtractor := func(permission *types.Permission) string { return permission.Name }

	utils.FromExtractor(pol.Resources, resourceExtractor).SortInPlace()

	for _, resource := range pol.Resources {
		utils.FromExtractor(resource.Relations, relationExtractor).SortInPlace()
		utils.FromExtractor(resource.Permissions, permissionExtractor).SortInPlace()
	}
	return specification.TransformerResult{
		Policy: pol,
	}, nil
}

func (t *SortTransformer) GetName() string { return "Sort Transformer" }
