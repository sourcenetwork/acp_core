package permission_parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermissionParser(t *testing.T) {
	perm := "owner + (that + resource->rel) - bar"

	tree, err := Parse(perm)

	want := "((owner + (that + resource->rel)) - bar)"
	require.NoError(t, err)
	require.Equal(t, want, tree.IntoPermissionExpr())
}
