package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResource_GetPermissionByName_ReturnsNilIfPermIsNil(t *testing.T) {
	resource := Resource{
		Permissions: nil,
	}

	perm := resource.GetPermissionByName("potato")

	require.Nil(t, perm)
}
