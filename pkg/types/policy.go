package types

import (
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

const OwnerRelationName = "owner"

// GetResourceByName returns a Resource with the given Name.
// If no resource is found with resourceName, return nil
func (pol *Policy) GetResourceByName(resourceName string) *Resource {
	for _, resource := range pol.Resources {
		if resource.Name == resourceName {
			return resource
		}
	}
	return nil
}

// ListResourcesNames returns a slice with the name of all defined resources
func (pol *Policy) ListResourcesNames() []string {
	return utils.MapSlice(pol.Resources, func(r *Resource) string { return r.Name })
}

// GetPermissionByName returns a Permission with `name`.
// If no Permission matches name, returns nil
func (res *Resource) GetPermissionByName(name string) *Permission {
	for _, permission := range res.Permissions {
		if permission.Name == name {
			return permission
		}
	}
	return nil
}

// GetRelationByName returns a Relation with `name`.
// If no Relation matches name, returns nil
func (res *Resource) GetRelationByName(name string) *Relation {
	if name == OwnerRelationName {
		return res.Owner
	}
	for _, relation := range res.Relations {
		if relation.Name == name {
			return relation
		}
	}
	return nil
}

// GetManagementRuleByName returns a ManagementPermission with `name`.
// If no ManagementPermission matches name, returns nil
func (res *Resource) GetManagementRuleByName(relation string) *ManagementRule {
	for _, rule := range res.ManagementRules {
		if rule.Relation == relation {
			return rule
		}
	}
	return nil
}

// ListRelationsNames returns the names of all Relations defined in the Resource
func (res *Resource) ListRelationsNames() []string {
	return utils.MapSlice(res.Relations, func(rel *Relation) string { return rel.Name })
}

// ListPermissionsNames returns the names of all Permissions defined in the Resource
func (res *Resource) ListPermissionsNames() []string {
	return utils.MapSlice(res.Permissions, func(rel *Permission) string { return rel.Name })
}
