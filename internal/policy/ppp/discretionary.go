package ppp

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	parser "github.com/sourcenetwork/acp_core/pkg/parser/permission_parser"
	"github.com/sourcenetwork/acp_core/pkg/transformer"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

var _ transformer.Transformer = (*DiscretionaryTransformer)(nil)

const (
	OwnerRelationName = "owner"
	OwnerDescription  = "owner relations represents the object owner"
)

var ErrDiscretionaryTransformer = errors.New("discretionary policy transformer", errors.ErrorType_BAD_INPUT)

// DiscretionaryTransformer is an essential part of acp_core which
// asures the owner relation exists for all resources in a policy.
// Further, it also guarantees that the owner authority over the actions
// which can be performed against an object are absolute.
//
// In practice it means that there exists a computed userset fetch rule
// for the owner at some point in the expression tree, adding it if necessary
type DiscretionaryTransformer struct{}

func (t *DiscretionaryTransformer) GetBaseError() error {
	return ErrDiscretionaryTransformer
}

func (t *DiscretionaryTransformer) Validate(policy types.Policy) *errors.MultiError {
	multiErr := errors.NewMultiError(ErrDiscretionaryTransformer)
	for _, resource := range policy.Resources {
		ownerRel := utils.FilterSlice(resource.Relations, func(r *types.Relation) bool { return r.Name == OwnerRelationName })
		if len(ownerRel) > 1 {
			err := fmt.Errorf("invalid policy: resource %v: multiple owner relations", resource.Name)
			multiErr.Append(err)
		}

		if len(ownerRel) == 0 {
			err := fmt.Errorf("invalid policy: resource %v: no owner relation", resource.Name)
			multiErr.Append(err)
		}
	}

	for _, resource := range policy.Resources {
		for _, permission := range resource.Permissions {
			tree, err := parser.Parse(permission.Expression)
			if err != nil {
				err := fmt.Errorf("parsing permission: resource %v: permission %v: %w", resource.Name, permission.Name, err)
				multiErr.Append(err)
				continue
			}

			if !checkOwnerIsTopNode(tree) {
				err := fmt.Errorf("invalid permission: resource %v: permission %v: expression does not contain owner as topmost allowed relation", resource.Name, permission.Name)
				multiErr.Append(err)
				continue
			}
		}
	}

	if len(multiErr.GetErrors()) > 0 {
		return multiErr
	}
	return nil
}

// Transform mutates all resources in a policy by asserting the owner relation exists, adding it otherwise.
// Futheremore, it modifies all permissions (if necessary),
// by adding the owner relation as one of the allowed relations.
func (t *DiscretionaryTransformer) Transform(policy types.Policy) (types.Policy, error) {
	// for all resources, add owner relation to it, if it doesn't exist
	for _, resource := range policy.Resources {
		ownerRel := utils.FilterSlice(resource.Relations, func(r *types.Relation) bool { return r.Name == OwnerRelationName })
		if len(ownerRel) > 1 {
			return types.Policy{}, errors.Wrap("invalid resource: multiple owner relations",
				ErrDiscretionaryTransformer,
				errors.Pair("resource", resource.Name))
		}
		if len(ownerRel) == 0 {
			rel := newOwnerRelation(&policy)
			ownerRel = append(ownerRel, rel)
			resource.Relations = append(resource.Relations, rel)
		}
	}

	// for all permissions in all resources
	// add a computed userset fetch rule as the toplevel operation
	for _, resource := range policy.Resources {
		for _, permission := range resource.Permissions {
			tree, err := parser.Parse(permission.Expression)
			if err != nil {
				return types.Policy{}, errors.Wrap("parsing permission", ErrDiscretionaryTransformer,
					errors.Pair("resource", resource.Name),
					errors.Pair("permission", permission.Name),
				)
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
func (t *DiscretionaryTransformer) transformFetchTree(tree *types.PermissionFetchTree) *types.PermissionFetchTree {
	if checkOwnerIsTopNode(tree) {
		return tree
	}

	return &types.PermissionFetchTree{
		Term: &types.PermissionFetchTree_CombNode{
			CombNode: &types.CombinationNode{
				Left:       buildFetchOwnerTree(),
				Combinator: types.Combinator_UNION,
				Right:      tree,
			},
		},
	}
}

// buildFetchOwnerTree returns a PermissionFetchTree with ComputedUserset owner as the single node
func buildFetchOwnerTree() *types.PermissionFetchTree {
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

func isComputedUsersetOwnerTree(tree *types.PermissionFetchTree) bool {
	if tree.GetOperation() == nil {
		return false
	}

	operation := tree.GetOperation()
	if operation.GetCu() == nil {
		return false
	}

	return operation.GetCu().Relation == OwnerRelationName
}

// checkOwnerIsTopNode performs a simples check which verifies
// whether got the given Fetch tree there is a computed userset rule
// for the owner relation as either a standalone node or
// as a top-level node for a union combination node.
//
// note this check is not comprehensive, as (owner + reader) + foo
// is still a valid permission expression, which would be rejected.
//
// TODO - add more compherensive check - either smt or an ad hoc algorithm
// which checks whether -owner exists at a higher level than an +owner
func checkOwnerIsTopNode(tree *types.PermissionFetchTree) bool {
	switch term := tree.Term.(type) {
	case *types.PermissionFetchTree_Operation:
		return isComputedUsersetOwnerTree(tree)
	case *types.PermissionFetchTree_CombNode:
		left := term.CombNode.Left
		right := term.CombNode.Right
		isUnionNode := term.CombNode.Combinator == types.Combinator_UNION
		isCUForOwner := (isComputedUsersetOwnerTree(left) || isComputedUsersetOwnerTree(right))
		return isCUForOwner && isUnionNode
	default:
		return false
	}
}
