package policy

import (
	"fmt"

	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/sourcenetwork/zanzi/pkg/domain"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

const defaultActorResourceName string = "actor"
const managementPermissionPrefix string = "_can_manage_"
const managementPermissionDoc string = "permission controls actors which are allowed to create relationships for the %v relation (permission was auto-generated by SourceHub)."

type factory struct{}

// NewPolicy creates a new policy from a marshal policy string.
// The policy is unmarshaled according to the given marshaling type and normalized.
func (f *factory) Create(policyIR PolicyIR, creator string, policyCounter uint64, creationTime *gogotypes.Timestamp) (*types.PolicyRecord, error) {

	policy := &types.Policy{
		Id:            "",
		Name:          policyIR.Name,
		Description:   policyIR.Description,
		CreationTime:  creationTime,
		Attributes:    policyIR.Attributes,
		Resources:     policyIR.Resources,
		ActorResource: policyIR.ActorResource,
		Creator:       creator,
	}
	f.normalize(policy)

	ider := policyIder{}
	policy.Id = ider.Id(policy, policyCounter)

	graph := buildManagementGraph(policy)
	f.registerOwnerAsManager(policy, graph)
	f.addManagementPermissions(policy, graph)

	return &types.PolicyRecord{
		Policy:          policy,
		ManagementGraph: graph,
	}, nil
}

// registerOwnerAsManager adds, for every Resource, every Relation as being managed by the "owner" Relation
func (f *factory) registerOwnerAsManager(policy *types.Policy, graph *types.ManagementGraph) {
	for _, resource := range policy.Resources {
		for _, relation := range resource.Relations {
			graph.RegisterManagedRel(resource.Name, OwnerRelation, relation.Name)
		}
	}
}

// creates a new permission for every relation in every resource.
// the management permission is used to verify which users can create relationships for a relation
func (f *factory) addManagementPermissions(policy *types.Policy, graph *types.ManagementGraph) {
	for _, resource := range policy.Resources {
		managementPermissions := make([]*types.Permission, 0, len(resource.Relations))
		for _, relation := range resource.Relations {
			permission := f.buildManagementPermission(resource.Name, relation, graph)
			managementPermissions = append(managementPermissions, permission)
		}
		resource.Permissions = append(resource.Permissions, managementPermissions...)
	}
}

func (f *factory) buildManagementPermission(resourceName string, relation *types.Relation, graph *types.ManagementGraph) *types.Permission {
	managerRelations := graph.GetManagers(resourceName, relation.Name)

	exprTree := f.buildRelationExpression(managerRelations)

	return &types.Permission{
		Name:       managementPermissionPrefix + relation.Name,
		Doc:        fmt.Sprintf(managementPermissionDoc, relation.Name),
		Expression: exprTree.RelationExpression(),
	}
}

func (f *factory) buildRelationExpression(relations []string) *domain.RelationExpressionTree {
	rel := relations[0]
	tree := &domain.RelationExpressionTree{
		Node: &domain.RelationExpressionTree_Rule{
			Rule: &domain.Rule{
				Rule: &domain.Rule_Cu{
					Cu: &domain.ComputedUserset{
						TargetRelation: rel,
					},
				},
			},
		},
	}
	relations = relations[1:]

	if len(relations) == 0 {
		return tree
	}

	return &domain.RelationExpressionTree{
		Node: &domain.RelationExpressionTree_OpNode{
			OpNode: &domain.OpNode{
				Left:     tree,
				Operator: domain.Operator_UNION,
				Right:    f.buildRelationExpression(relations),
			},
		},
	}
}

// normalize normalizes a policy by setting default values for optional fields.
func (f *factory) normalize(pol *types.Policy) {
	if pol.ActorResource == nil {
		pol.ActorResource = &types.ActorResource{
			Name: defaultActorResourceName,
		}
	}

	// policy is sorted before building id to ensure determinism
	pol.Sort()
}
