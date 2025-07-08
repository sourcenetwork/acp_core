package ppp

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/specification"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
	"k8s.io/apimachinery/pkg/util/sets"
)

var _ specification.Transformer = (*DecentralizedAdminTransformer)(nil)

const (
	managementPermissionPrefix string = "_can_manage_"
	managementPermissionDoc    string = "permission controls actors which are allowed to create relationships for the %v relation (permission was auto-generated)."
)

var ErrAdministrationTransformer = errors.New("decentralized administration transformer", errors.ErrorType_BAD_INPUT)

// DecentralizedAdminTransformer transforms a Policy by adding the management permissions,
// which controls which actors are able to mutate relationships to an object
type DecentralizedAdminTransformer struct{}

// Validate asserts that for all relations in a Policy,
// there exists a corresponding management permission
func (t *DecentralizedAdminTransformer) Validate(policy types.Policy) *errors.MultiError {
	multiErr := errors.NewMultiError(ErrAdministrationTransformer)

	for _, resource := range policy.Resources {
		permissionNames := utils.MapSlice(resource.Permissions, func(p *types.Permission) string {
			return p.Name
		})
		permissionSet := sets.New(permissionNames...)

		for _, relation := range resource.Relations {
			managementPermName := t.buildManagementPermissionName(relation.Name)
			if !permissionSet.Has(managementPermName) {
				err := fmt.Errorf("management permission not found: resource %v: relation %v", resource.Name, relation.Name)
				multiErr.Append(err)
			}
		}
	}

	if len(multiErr.GetErrors()) > 0 {
		return multiErr
	}

	return nil
}

// Transform processes the `manages` directives in a Policy
// and creates the required management permissions for it
func (t *DecentralizedAdminTransformer) Transform(pol types.Policy) (types.Policy, error) {
	graph := &types.ManagementGraph{}
	graph.LoadFromPolicy(&pol)
	err := graph.IsWellFormed()
	if err != nil {
		return types.Policy{}, errors.Wrap("invalid manages definition: "+err.Error(), ErrAdministrationTransformer)
	}

	for _, resource := range pol.Resources {
		for _, relation := range resource.Relations {
			managementPermission := t.buildManagementPermission(resource.Name, relation, graph)
			resource.Permissions = append(resource.Permissions, managementPermission)
		}
	}

	return pol, nil
}

func (t *DecentralizedAdminTransformer) buildManagementPermission(resourceName string, relation *types.Relation, graph *types.ManagementGraph) *types.Permission {
	managerRelations := graph.GetManagers(resourceName, relation.Name)

	exprTree := t.buildRelationExpression(managerRelations)

	return &types.Permission{
		Name:       t.buildManagementPermissionName(relation.Name),
		Doc:        fmt.Sprintf(managementPermissionDoc, relation.Name),
		Expression: exprTree.IntoPermissionExpr(),
	}
}

func (t *DecentralizedAdminTransformer) buildRelationExpression(relations []string) *types.PermissionFetchTree {
	if len(relations) == 0 {
		return newFetchOwnerTree()
	}

	tree := &types.PermissionFetchTree{
		Term: &types.PermissionFetchTree_Operation{
			Operation: &types.FetchOperation{
				Operation: &types.FetchOperation_Cu{
					Cu: &types.ComputedUsersetNode{
						Relation: relations[0],
					},
				},
			},
		},
	}
	for _, relation := range relations[1:len(relations)] {
		node := &types.PermissionFetchTree{
			Term: &types.PermissionFetchTree_Operation{
				Operation: &types.FetchOperation{
					Operation: &types.FetchOperation_Cu{
						Cu: &types.ComputedUsersetNode{
							Relation: relation,
						},
					},
				},
			},
		}
		tree = &types.PermissionFetchTree{
			Term: &types.PermissionFetchTree_CombNode{
				CombNode: &types.CombinationNode{
					Left:       tree,
					Combinator: types.Combinator_UNION,
					Right:      node,
				},
			},
		}
	}

	return &types.PermissionFetchTree{
		Term: &types.PermissionFetchTree_CombNode{
			CombNode: &types.CombinationNode{
				Left:       tree,
				Combinator: types.Combinator_UNION,
				Right:      newFetchOwnerTree(),
			},
		},
	}
}

func (t *DecentralizedAdminTransformer) buildManagementPermissionName(relationName string) string {
	return managementPermissionPrefix + relationName
}

func (t *DecentralizedAdminTransformer) GetBaseError() error { return ErrAdministrationTransformer }
