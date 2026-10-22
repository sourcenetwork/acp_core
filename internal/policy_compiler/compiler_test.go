package policy_compiler

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

func compile(t *testing.T, policy string) (*Model, error) {
	t.Helper()
	return Compile(policy, types.PolicyMarshalingType_YAML)
}

func Test_Compile_RecursiveTTUBecomesClosure(t *testing.T) {
	model, err := compile(t, `
name: test
resources:
- name: folder
  relations:
  - name: parent
    types: [folder]
  - name: reader
    types: [actor]
  permissions:
  - name: read
    expr: reader + parent->read
`)
	require.NoError(t, err)

	read := model.Resources[0].Permissions[0]
	closure, ok := read.Expr.(*Closure)
	require.True(t, ok, "expected closure, got %T", read.Expr)
	require.Len(t, closure.Fields, 1)
	require.Equal(t, "parent", closure.Fields[0].Name)
}

func Test_Compile_NestedUsersetBecomesClosure(t *testing.T) {
	model, err := compile(t, `
name: test
resources:
- name: group
  relations:
  - name: member
    types: [actor, group->member]
`)
	require.NoError(t, err)

	member := model.Resources[0].Relations[1]
	require.Equal(t, "member", member.Name)
	_, ok := member.Expr.(*Closure)
	require.True(t, ok, "expected closure, got %T", member.Expr)
}

func Test_Compile_FieldsAreNamedPerRestriction(t *testing.T) {
	model, err := compile(t, `
name: test
resources:
- name: group
  relations:
  - name: member
    types: [actor]
- name: file
  relations:
  - name: parent
    types: [file]
  - name: reader
    types: [actor, group->member, group]
`)
	require.NoError(t, err)

	var names []string
	for _, f := range model.Resources[0].Fields {
		names = append(names, f.Name)
	}
	require.Equal(t, []string{"owner", "parent", "reader", "reader_group_member", "reader_group"}, names)
}

func Test_Compile_TTUWithMultipleTypesErrors(t *testing.T) {
	_, err := compile(t, `
name: test
resources:
- name: folder
  relations:
  - name: reader
    types: [actor]
- name: file
  relations:
  - name: parent
    types: [folder, file]
  permissions:
  - name: read
    expr: parent->reader
`)
	require.ErrorIs(t, err, ErrPolicyCompiler)
	require.ErrorContains(t, err, "relation parent must have exactly one resource type restriction")
}

func Test_Compile_TTUOverUsersetErrors(t *testing.T) {
	_, err := compile(t, `
name: test
resources:
- name: group
  relations:
  - name: member
    types: [actor]
- name: file
  relations:
  - name: team
    types: [group->member]
  permissions:
  - name: read
    expr: team->member
`)
	require.ErrorIs(t, err, ErrPolicyCompiler)
	require.ErrorContains(t, err, "relation team must have exactly one resource type restriction")
}

func Test_Compile_RecursionUnderDifferenceErrors(t *testing.T) {
	_, err := compile(t, `
name: test
resources:
- name: folder
  relations:
  - name: parent
    types: [folder]
  - name: reader
    types: [actor]
  - name: banned
    types: [actor]
  permissions:
  - name: read
    expr: (reader + parent->read) - banned
`)
	require.ErrorIs(t, err, ErrPolicyCompiler)
	require.ErrorContains(t, err, "recursive definition not supported: folder.read -> folder.read")
}

func Test_Compile_MutualRecursionErrors(t *testing.T) {
	_, err := compile(t, `
name: test
resources:
- name: folder
  relations:
  - name: file
    types: [file]
  permissions:
  - name: read
    expr: file->read
- name: file
  relations:
  - name: folder
    types: [folder]
  permissions:
  - name: read
    expr: folder->read
`)
	require.ErrorIs(t, err, ErrPolicyCompiler)
	require.ErrorContains(t, err, "recursive definition not supported")
}
