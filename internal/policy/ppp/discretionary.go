package ppp

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/specification"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	parser "github.com/sourcenetwork/acp_core/pkg/parser/permission_parser"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

var _ specification.Transformer = (*DiscretionaryTransformer)(nil)

func GetOwnerRelation(actorResourceName string, managedRels []string) *types.Relation {
	return &types.Relation{
		Name: OwnerRelationName,
		Doc:  ownerDescription,
		VrTypes: []*types.Restriction{
			{
				ResourceName: actorResourceName,
			},
		},
		Manages: managedRels,
	}
}

const (
	OwnerRelationName = "owner"
	ownerDescription  = "owner relations represents the object owner"
)

var ErrDiscretionaryTransformer = errors.New("discretionary policy transformer", errors.ErrorType_BAD_INPUT)
var ErrResourceContainsOwner = errors.New("invalid resource: resource includes reserved `owner` relation: rename `owner` to another relation name", errors.ErrorType_BAD_INPUT)
var ErrPermissionReferencesOwner = errors.New("invalid permission: permission cannot reference `owner` relation as part of the computed expression. Try removing `owner` or asserting the tuple to userset expression is correct", errors.ErrorType_BAD_INPUT)

// DiscretionaryTransformer is an essential part of acp_core which
// asures the owner relation exists for all resources in a policy.
// Further, it also guarantees that the owner authority over the actions
// which can be performed against an object are absolute.
//
// In practice it means that there exists a computed userset fetch rule
// for the owner at some point in the expression tree, adding it if necessary
type DiscretionaryTransformer struct{}

func (t *DiscretionaryTransformer) GetName() string {
	return "Discretionary transformer"
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

// findOwnerViolationsInParseTree returns a term in a PermissionFetchTree
// which references owner in a way which violates the discretionry transformer.
// ie. only owners in another resource (tuple to userset) can be referenced in a permission expr.
func findOwnerViolationsInParseTree(tree *types.PermissionFetchTree) error {
	switch term := tree.Term.(type) {
	case *types.PermissionFetchTree_CombNode:
		err := findOwnerViolationsInParseTree(term.CombNode.Left)
		if err != nil {
			return err
		}
		return findOwnerViolationsInParseTree(term.CombNode.Right)
	case *types.PermissionFetchTree_Operation:
		switch op := term.Operation.Operation.(type) {
		case *types.FetchOperation_Cu:
			if op.Cu.Relation == OwnerRelationName {
				return fmt.Errorf("violation found: cannot reference `owner` as part of expression")
			}
		case *types.FetchOperation_Ttu:
			// we only care about the first relation in a TTU,
			// because it's perfectly valid to reference the owner of another
			// resource as the inheritance target
			if op.Ttu.LookupRelation == OwnerRelationName {
				return fmt.Errorf("violation found: `owner` relation cannot be used as the lookup relation in a tuple to userset")
			}
			// other cases are uninteresting
		}
	}
	return nil
}

func (t *DiscretionaryTransformer) validatePermissionDoesNotReferenceOwner(expr string) error {
	tree, err := parser.Parse(expr)
	if err != nil {
		return fmt.Errorf("invalid permission: parsing error: %v", err)
	}
	err = findOwnerViolationsInParseTree(tree)
	if err != nil {
		return fmt.Errorf("invalid permission: %w", err)
	}
	return nil
}

// Transform mutates all resources in a policy by asserting the owner relation exists, adding it otherwise.
// Futheremore, it modifies all permissions (if necessary),
// by adding the owner relation as one of the allowed relations.
func (t *DiscretionaryTransformer) Transform(policy types.Policy) (specification.TransformerResult, error) {
	res := specification.TransformerResult{}

	// pre-validation: assert that resoruces don not include an `owner` relation
	for _, resource := range policy.Resources {
		rel := resource.GetRelationByName(OwnerRelationName)
		if rel != nil {
			return res, errors.Attrs(ErrResourceContainsOwner,
				errors.Pair("resource", resource.Name))
		}
	}

	// add special owner relation to every resource
	for _, resource := range policy.Resources {
		relNames := utils.MapSlice(resource.Relations, func(r *types.Relation) string { return r.Name })
		resource.Owner = GetOwnerRelation(policy.ActorResource.Name, relNames)
	}

	// walk through permissions asserting it does not reference owner
	for _, resource := range policy.Resources {
		for _, permission := range resource.Permissions {
			err := t.validatePermissionDoesNotReferenceOwner(permission.Expression)
			if err != nil {
				return res, errors.Attrs(ErrPermissionReferencesOwner,
					errors.Pair("resource", resource.Name),
					errors.Pair("permission", permission.Name))
			}
		}
	}

	// for all permissions in all resources
	// add a computed userset fetch rule as the toplevel operation
	for _, resource := range policy.Resources {
		for _, permission := range resource.Permissions {
			tree, err := parser.Parse(permission.Expression)
			if err != nil {
				return res, errors.Wrap("parsing permission", errors.Wrap(err.Error(), ErrDiscretionaryTransformer,
					errors.Pair("resource", resource.Name),
					errors.Pair("permission", permission.Name),
				))
			}
			tree, modified := t.transformFetchTree(tree)
			expr := tree.IntoPermissionExpr()
			permission.EffectiveExpression = expr
			if modified {
				msg := fmt.Sprintf("added owner relation to permission: resource %v: permission %v",
					resource.Name, permission.Name)
				res.Messages = append(res.Messages, msg)
			}
		}
	}

	res.Policy = policy
	return res, nil
}

// transformFetchTree adds computed userset instruction for owner as leftmost node
// if the tree already meets this criteria, this is a noop
//
// this transformation is very primitive, as there are several
// relations which would still be valid which aren't accounted for (eg. owner + foo + barr),
// however it is the simplest way to ensure owners have full access
// as there are several subtle expressions which could remove owner access
// eg. owner - (something & owner)
func (t *DiscretionaryTransformer) transformFetchTree(tree *types.PermissionFetchTree) (_ *types.PermissionFetchTree, modified bool) {
	if checkOwnerIsTopNode(tree) {
		return tree, false
	}

	return &types.PermissionFetchTree{
		Term: &types.PermissionFetchTree_CombNode{
			CombNode: &types.CombinationNode{
				Left:       newFetchOwnerTree(),
				Combinator: types.Combinator_UNION,
				Right:      tree,
			},
		},
	}, true
}

// newFetchOwnerTree returns a PermissionFetchTree with ComputedUserset owner as the single node
func newFetchOwnerTree() *types.PermissionFetchTree {
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

// isComputedUsersetOwnerTree verifies whether the given fetch tree
// is a single node tree for a computed userset operation for the owner relation
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
