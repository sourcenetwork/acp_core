package specification

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	parser "github.com/sourcenetwork/acp_core/pkg/parser/permission_parser"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
	"k8s.io/apimachinery/pkg/util/sets"
)

var _ Requirement = (*defaultPermissionsRequirement)(nil)
var _ Transformer = (*writeImpliesReadTransform)(nil)

// NewDefraSpecification returns an instace of Specificatoin
// to validate Defra compliant Policies
func NewDefraSpecification() Specification {
	return newSpecification(types.PolicySpecificationType_DEFRA_SPEC,
		[]Requirement{
			&defaultPermissionsRequirement{},
		},
		[]Transformer{
			&writeImpliesReadTransform{},
		},
	)
}

const (
	DefraReadPermissionName  = "read"
	DefraWritePermissionName = "write"
)

// ErrDefraSpec is the base error for the DefraSpec implementation
var ErrDefraSpec = errors.New("defra policy specification", errors.ErrorType_BAD_INPUT)

// RequiredPermissiosn are the set of permissions which all resources
// must include in a Defra compliant Policy.
var RequiredPermissions = []string{DefraReadPermissionName, DefraWritePermissionName}

// NewDefraSpec returns an instance of DefraSpec
func NewDefraSpec() Requirement {
	return &defaultPermissionsRequirement{}
}

// defaultPermissionsRequirement implements the Specification interface for Defra compliant Policies
//
// Defra compliant Policies require that all resources contain a set of pre-determined
// permissions.
type defaultPermissionsRequirement struct{}

func (s *defaultPermissionsRequirement) Validate(pol types.Policy) *errors.MultiError {
	multiErr := errors.NewMultiError(ErrDefraSpec)

	for _, resource := range pol.Resources {
		permissions := utils.MapSlice(resource.Permissions, func(p *types.Permission) string { return p.Name })
		permissionsSet := sets.New(permissions...)
		for _, permission := range RequiredPermissions {
			if !permissionsSet.Has(permission) {
				violation := fmt.Errorf("resource %v: missing required permission: %v", resource.Name, permission)
				multiErr.Append(violation)
			}
		}
	}
	if len(multiErr.GetErrors()) != 0 {
		return multiErr
	}
	return nil
}

func (s *defaultPermissionsRequirement) GetName() string { return "Defra Spec" }

type writeImpliesReadTransform struct {
	result TransformerResult
}

func (s *writeImpliesReadTransform) GetName() string {
	return "Defra Spec"
}

func (s *writeImpliesReadTransform) Validate(pol types.Policy) *errors.MultiError {
	return nil
}

func (s *writeImpliesReadTransform) Transform(pol types.Policy) (TransformerResult, error) {
	result := TransformerResult{}
	for _, resource := range pol.Resources {
		for _, permission := range resource.Permissions {
			if permission.Name != DefraReadPermissionName {
				continue
			}
			tree, err := parser.Parse(permission.Expression)
			if err != nil {
				return result, errors.Wrap("parsing permission", ErrDefraSpec,
					errors.Pair("resource", resource.Name),
					errors.Pair("permission", permission.Name),
				)
			}
			tree = s.transformTree(tree)
			permission.Expression = tree.IntoPermissionExpr()
		}
	}
	result.Policy = pol
	return result, nil
}

func (s *writeImpliesReadTransform) transformTree(tree *types.PermissionFetchTree) *types.PermissionFetchTree {
	// defra transformers should run after the discretionary ones
	// meaning the top level fetch instruction should be owner computed userset
	if s.isCUOwnerTree(tree) {
		return &types.PermissionFetchTree{
			Term: &types.PermissionFetchTree_CombNode{
				CombNode: &types.CombinationNode{
					Left:       tree,
					Combinator: types.Combinator_UNION,
					Right:      s.writeCUNode(),
				},
			},
		}
	}

	comb := tree.GetCombNode()
	ownerTree, remainder := comb.Left, comb.Right
	if s.isCUOwnerTree(comb.Right) {
		ownerTree = comb.Right
		remainder = comb.Left
	}

	remainder = &types.PermissionFetchTree{
		Term: &types.PermissionFetchTree_CombNode{
			CombNode: &types.CombinationNode{
				Left:       s.writeCUNode(),
				Combinator: types.Combinator_UNION,
				Right:      remainder,
			},
		},
	}

	return &types.PermissionFetchTree{
		Term: &types.PermissionFetchTree_CombNode{
			CombNode: &types.CombinationNode{
				Left:       ownerTree,
				Combinator: types.Combinator_UNION,
				Right:      remainder,
			},
		},
	}
}

func (s *writeImpliesReadTransform) writeCUNode() *types.PermissionFetchTree {
	return &types.PermissionFetchTree{
		Term: &types.PermissionFetchTree_Operation{
			Operation: &types.FetchOperation{
				Operation: &types.FetchOperation_Cu{
					Cu: &types.ComputedUsersetNode{
						Relation: DefraWritePermissionName,
					},
				},
			},
		},
	}
}

func (s *writeImpliesReadTransform) isCUOwnerTree(tree *types.PermissionFetchTree) bool {
	if tree.GetOperation() == nil {
		return false
	}
	op := tree.GetOperation()
	if op.GetCu() == nil {
		return false
	}
	cu := op.GetCu()
	return cu.Relation == "owner"
}
