package ppp

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
	"k8s.io/apimachinery/pkg/util/sets"
)

const managementPermissionPrefix string = "_can_manage_"
const managementPermissionDoc string = "permission controls actors which are allowed to create relationships for the %v relation (permission was auto-generated)."

var _ Transformer = (*DecentralizedAdminTransformer)(nil)

type DecentralizedAdminTransformer struct{}

func (t *DecentralizedAdminTransformer) Name() string {
	return "Decentralized Administration"
}

func (t *DecentralizedAdminTransformer) Validate(policy *types.Policy) *errors.MultiError {
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

func (t *DecentralizedAdminTransformer) Transform(provider PolicyProvider) (*types.Policy, *errors.MultiError) {
	pol := provider()
	graph := &types.ManagementGraph{}
	graph.LoadFromPolicy(pol)
	err := graph.IsWellFormed()
	if err != nil {
		return nil, errors.NewMultiError(ErrAdministrationTransformer, err)
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

// creates a new permission for every relation in every resource.
// the management permission is used to verify which users can create relationships for a relation
func (t *DecentralizedAdminTransformer) addManagementPermissions(policy *types.Policy, graph *types.ManagementGraph) {
	for _, resource := range policy.Resources {
		managementPermissions := make([]*types.Permission, 0, len(resource.Relations))
		for _, relation := range resource.Relations {
			permission := t.buildManagementPermission(resource.Name, relation, graph)
			managementPermissions = append(managementPermissions, permission)
		}
		resource.Permissions = append(resource.Permissions, managementPermissions...)
	}
}

func (t *DecentralizedAdminTransformer) buildRelationExpression(relations []string) *types.PermissionFetchTree {
	tree := &types.PermissionFetchTree{
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

	for _, relation := range relations {
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

	return tree
}

func (t *DecentralizedAdminTransformer) buildManagementPermissionName(relationName string) string {
	return managementPermissionPrefix + relationName
}
