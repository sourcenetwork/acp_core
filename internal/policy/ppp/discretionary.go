package ppp

import (
	"fmt"

	parser "github.com/sourcenetwork/acp_core/internal/parser/permission_parser"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

const OwnerRelationName = "owner"
const OwnerDescription = "owner relations represents the object owner"

var _ Transformer = (*DiscretionaryTransformer)(nil)

type DiscretionaryTransformer struct{}

func (t *DiscretionaryTransformer) Name() string {
	return "Discretionary Transformer"
}

func (t *DiscretionaryTransformer) Validate(policy *types.Policy) []error {
	var violations []error
	for _, resource := range policy.Resources {
		ownerRel := utils.FilterSlice(resource.Relations, func(r *types.Relation) bool { return r.Name == OwnerRelationName })
		if len(ownerRel) > 1 {
			err := fmt.Errorf("invalid policy: resource %v: multiple owner relations", resource.Name)
			violations = append(violations, err)
		}

		if len(ownerRel) == 0 {
			err := fmt.Errorf("invalid policy: resource %v: no owner relation")
			violations = append(violations, err)
		}
	}

	for _, resource := range policy.Resources {
		for _, permission := range resource.Permissions {
			tree, report := parser.Parse(permission.Expression)
			if report.HasError() {
				err := fmt.Errorf("parsing permission: resource %v: permission %v: %w", resource.Name, permission.Name, report.Error())
				violations = append(violations, err)
				continue
			}

			if !t.checkOwnerIsAllowed(tree) {
				err := fmt.Errorf("invalid permission: resource %v: permission %v: expression does not contain owner as topmost allowed relation", resource.Name, permission.Name)
				violations = append(violations, err)
				continue
			}
		}
	}

	return violations
}

func (t *DiscretionaryTransformer) Transform(provider PolicyProvider) (*types.Policy, error) {
	policy := provider()

	// for all resources, add owner relation to it, if it doesn't exist
	for _, resource := range policy.Resources {
		ownerRel := utils.FilterSlice(resource.Relations, func(r *types.Relation) bool { return r.Name == OwnerRelationName })
		if len(ownerRel) > 1 {
			return nil, fmt.Errorf("invalid policy: resource %v: multiple owner relations", resource.Name)
		}
		if len(ownerRel) == 0 {
			rel := newOwnerRelation(policy)
			ownerRel = append(ownerRel, rel)
			resource.Relations = append(resource.Relations, rel)
		}
	}

	// for all permissions in all resources, add owner as the leftmost term
	for _, resource := range policy.Resources {
		for _, permission := range resource.Permissions {
			tree, report := parser.Parse(permission.Expression)
			if report.HasError() {
				return nil, fmt.Errorf("parsing permission: resource %v: permission %v: %w", resource.Name, permission.Name, report.Error())
			}
			tree = t.transformFetchTree(tree)
			expr := tree.IntoPermissionExpr()
			permission.Expression = expr
		}
	}

	return policy, nil
}

// transformFetchTree adds computed userset instruction for owner as leftmost node
// if the tree already meets this criteria, this is a noop
func (t *DiscretionaryTransformer) checkOwnerIsAllowed(tree *types.PermissionFetchTree) bool {
	var node, parent *types.PermissionFetchTree = tree, nil
	for {
		if node.GetOperation() != nil {
			break
		} else {
			comb := node.GetCombNode()
			parent = node
			node = comb.Left
		}
	}

	isNodeCUOwner := node.GetOperation().GetCu() != nil && node.GetOperation().GetCu().Relation == OwnerRelationName
	isParentNilOrOwnerIsAddedToParentCombNode := parent == nil || parent.GetCombNode().Combinator == types.Combinator_UNION

	// owner is already included
	return isNodeCUOwner && isParentNilOrOwnerIsAddedToParentCombNode
}

// transformFetchTree adds computed userset instruction for owner as leftmost node
// if the tree already meets this criteria, this is a noop
func (t *DiscretionaryTransformer) transformFetchTree(tree *types.PermissionFetchTree) *types.PermissionFetchTree {
	var node, parent *types.PermissionFetchTree = tree, nil
	for {
		if node.GetOperation() != nil {
			break
		} else {
			comb := node.GetCombNode()
			parent = node
			node = comb.Left
		}
	}

	isNodeCUOwner := node.GetOperation().GetCu() != nil && node.GetOperation().GetCu().Relation == OwnerRelationName
	isParentNilOrOwnerIsAddedToParentCombNode := parent == nil || parent.GetCombNode().Combinator == types.Combinator_UNION

	// noop, owner is already included
	if isNodeCUOwner && isParentNilOrOwnerIsAddedToParentCombNode {
		return tree
	}

	combNode := &types.PermissionFetchTree{
		Term: &types.PermissionFetchTree_CombNode{
			CombNode: &types.CombinationNode{
				Left:       t.buildFetchOwnerTree(),
				Combinator: types.Combinator_UNION,
				Right:      node,
			},
		},
	}
	if parent == nil {
		return combNode
	}
	parent.GetCombNode().Left = combNode
	return tree
}

// buildFetchOwnerTree returns a PermissionFetchTree with ComputedUserset owner as the single node
func (t *DiscretionaryTransformer) buildFetchOwnerTree() *types.PermissionFetchTree {
	return &types.PermissionFetchTree{
		Term: &types.PermissionFetchTree_Operation{
			Operation: &types.FetchOperation{
				Operation: &types.FetchOperation_Cu{
					Cu: &types.ComputedUsersetNode{
						Relation: OwnerRelationName,
					},
				},
			},
		},
	}
}

// newOwnerRelation returns a default relation for the resource owner
func newOwnerRelation(policy *types.Policy) *types.Relation {
	return &types.Relation{
		Name: OwnerRelationName,
		Doc:  OwnerDescription,
		VrTypes: []*types.Restriction{
			{
				ResourceName: policy.ActorResource.Name,
			},
		},
	}
}
