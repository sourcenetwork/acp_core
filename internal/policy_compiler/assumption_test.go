package policy_compiler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const assumptionsPolicy = `
name: test
resources:
- name: group
  relations:
  - name: member
    types: [actor, group->member]
- name: folder
  relations:
  - name: parent
    types: [folder]
- name: doc
  relations:
  - name: folder
    types: [folder]
  - name: editor
    types: [actor, group->member]
  - name: viewer
    types: [actor]
`

func lowerAssumptions(t *testing.T, definition string) ([]Assumption, error) {
	t.Helper()
	model, err := compile(t, assumptionsPolicy)
	require.NoError(t, err)
	assumptions, err := UnmarshalAssumptions(definition)
	if err != nil {
		return nil, err
	}
	return LowerAssumptions(model, assumptions)
}

func Test_Assumptions_Lower(t *testing.T) {
	lowered, err := lowerAssumptions(t, `
assumptions:
- disjoint: [doc.owner, doc.editor, doc.viewer]
- at_most_one: doc.folder
- acyclic: folder.parent
- acyclic: group.member
`)
	require.NoError(t, err)
	require.Len(t, lowered, 4)

	disjoint := lowered[0].(*Disjoint)
	var pairs []string
	for _, p := range disjoint.Pairs {
		pairs = append(pairs, p[0].Name+"/"+p[1].Name)
	}
	// only fields holding the same kind of subject are paired
	require.Equal(t, []string{"owner/editor", "owner/viewer", "editor/viewer"}, pairs)

	require.Len(t, lowered[1].(*AtMostOne).Fields, 1)
	require.Equal(t, "parent", lowered[2].(*Acyclic).Fields[0].Name)
	require.Equal(t, "member_group_member", lowered[3].(*Acyclic).Fields[0].Name)
}

func Test_Assumptions_UnmarshalErrors(t *testing.T) {
	cases := map[string]string{
		"several kinds in one entry": `
assumptions:
- disjoint: [doc.owner, doc.editor]
  acyclic: folder.parent`,
		"invalid relation reference": `
assumptions:
- at_most_one: doc`,
		"disjoint across resources": `
assumptions:
- disjoint: [doc.owner, folder.owner]`,
		"unknown key": `
assumptions:
- unique: doc.folder`,
	}
	for name, definition := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := UnmarshalAssumptions(definition)
			require.Error(t, err)
		})
	}
}

func Test_Assumptions_LowerErrors(t *testing.T) {
	cases := map[string]struct{ definition, msg string }{
		"undefined relation": {`
assumptions:
- at_most_one: doc.parent`, "undefined relation doc.parent"},
		"undefined resource": {`
assumptions:
- acyclic: file.parent`, "undefined resource file"},
		"relations never share subjects": {`
assumptions:
- disjoint: [doc.folder, doc.viewer]`, "relations folder and viewer never share a subject type"},
		"acyclic relation to another resource": {`
assumptions:
- acyclic: doc.folder`, "acyclic doc.folder: relation never points at a doc"},
		"single relation disjoint": {`
assumptions:
- disjoint: [doc.viewer]`, "expected at least two relations"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := lowerAssumptions(t, c.definition)
			require.ErrorIs(t, err, ErrPolicyCompiler)
			require.ErrorContains(t, err, c.msg)
		})
	}
}
