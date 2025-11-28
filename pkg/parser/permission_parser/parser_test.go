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

func TestPermission_InvalidSymbol(t *testing.T) {
	perm := "abc ^ something"

	tree, err := Parse(perm)

	require.Error(t, err)
	require.Nil(t, tree)
}

func TestPermission_SingleLetterRelation(t *testing.T) {
	perm := "w"

	tree, err := Parse(perm)

	want := "w"
	require.NoError(t, err)
	require.Equal(t, want, tree.IntoPermissionExpr())
}

func TestPermission_EmptyProduction_ReturnsNilTree(t *testing.T) {
	perm := ""

	tree, err := Parse(perm)

	require.NoError(t, err)
	require.Nil(t, tree.Term)
}
