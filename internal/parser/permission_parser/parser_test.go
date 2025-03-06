package permission_parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Maybe use a LISP-like / Reverse Polish Notation tree representation to ease asserting
func TestPermissionParser(t *testing.T) {
	perm := "owner + (that + resource->rel) - bar"

	tree, report := Parse(perm)

	want := "((owner + (that + resource->rel)) - bar)"
	require.False(t, report.HasError())
	require.Equal(t, want, tree.IntoPermissionExpr())
}
