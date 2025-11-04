package types

import (
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

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
	for _, relation := range res.Relations {
		if relation.Name == name {
			return relation
		}
	}
	return nil
}

// GetManagementPermissionByName returns a ManagementPermission with `name`.
// If no ManagementPermission matches name, returns nil
func (res *Resource) GetManagementPermissionByName(name string) *ManagementPermission {
	for _, perm := range res.ManagementPermissions {
		if perm.Name == name {
			return perm
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
