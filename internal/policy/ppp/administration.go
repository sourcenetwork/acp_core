package ppp

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/specification"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ specification.Transformer = (*DecentralizedAdminTransformer)(nil)

var ErrAdministrationTransformer = errors.New("decentralized administration transformer", errors.ErrorType_BAD_INPUT)

// DecentralizedAdminTransformer transforms a Policy by adding the management permissions,
// which controls which actors are able to mutate relationships to an object
type DecentralizedAdminTransformer struct{}

// Validate asserts that for all relations in a Policy,
// there exists a corresponding management permission
func (t *DecentralizedAdminTransformer) Validate(policy types.Policy) *errors.MultiError {
	multiErr := errors.NewMultiError(ErrAdministrationTransformer)

	for _, resource := range policy.Resources {
		mgmtPermissions := make(map[string]struct{})
		for _, perm := range resource.ManagementRules {
			mgmtPermissions[perm.Name] = struct{}{}
		}

		for _, rel := range resource.Relations {
			_, ok := mgmtPermissions[rel.Name]
			if !ok {
				err := fmt.Errorf("management permission not found: resource %v: relation %v", resource.Name, rel.Name)
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
func (t *DecentralizedAdminTransformer) Transform(pol types.Policy) (specification.TransformerResult, error) {
	res := specification.TransformerResult{}
	graph := &types.ManagementGraph{}
	graph.LoadFromPolicy(&pol)
	err := graph.IsWellFormed()
	if err != nil {
		return res, errors.Wrap("invalid manages definition: "+err.Error(), ErrAdministrationTransformer)
	}

	for _, resource := range pol.Resources {
		perms := make([]*types.ManagementRule, 0, len(resource.Relations))
		for _, relation := range resource.Relations {
			perm := t.buildManagementPermission(resource.Name, relation, graph)
			perms = append(perms, perm)
		}
		resource.ManagementRules = perms
	}

	res.Policy = pol
	return res, nil
}

func (t *DecentralizedAdminTransformer) buildManagementPermission(resourceName string, relation *types.Relation, graph *types.ManagementGraph) *types.ManagementRule {
	managerRelations := graph.GetManagers(resourceName, relation.Name)

	exprTree := t.buildRelationExpression(managerRelations)

	return &types.ManagementRule{
		Name:       relation.Name,
		Expression: exprTree.IntoPermissionExpr(),
		Relations:  managerRelations,
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

func (t *DecentralizedAdminTransformer) GetName() string { return "Decentralized Administrator" }
