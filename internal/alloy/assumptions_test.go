package alloy

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	pc "github.com/sourcenetwork/acp_core/internal/policy_compiler"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

const workspaceAssumptions = `
assumptions:
- disjoint: [doc.owner, doc.reviewer, doc.banned]
- disjoint: [group.admin, group.member]
- at_most_one: doc.folder
- acyclic: folder.parent
- acyclic: group.member
`

func emitWorkspaceAssumptions(t *testing.T) (*pc.Model, File) {
	t.Helper()
	model, err := pc.Compile(workspacePolicy, types.PolicyMarshalingType_YAML)
	require.NoError(t, err)
	assumptions, err := pc.UnmarshalAssumptions(workspaceAssumptions)
	require.NoError(t, err)
	lowered, err := pc.LowerAssumptions(model, assumptions)
	require.NoError(t, err)
	return model, EmitAssumptions(model, "", lowered)
}

func Test_EmitAssumptions_MatchesGoldenFile(t *testing.T) {
	_, f := emitWorkspaceAssumptions(t)

	path := filepath.Join("testdata", f.Name)
	if *update {
		require.NoError(t, os.WriteFile(path, []byte(f.Content), 0644))
	}
	golden, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, string(golden), f.Content, "%v differs from golden file, run with -update to regenerate", f.Name)
}

// Test_Assumptions_RestrictInstances asserts that instances violating an assumption
// exist without the assumptions and can't be generated with them
func Test_Assumptions_RestrictInstances(t *testing.T) {
	if _, err := exec.LookPath("alloy"); err != nil {
		t.Skip("alloy not found in PATH")
	}
	model, assumptions := emitWorkspaceAssumptions(t)
	dir := t.TempDir()
	for _, f := range append(Emit(model, ""), assumptions) {
		require.NoError(t, os.WriteFile(filepath.Join(dir, f.Name), []byte(f.Content), 0644))
	}

	violations := map[string]string{
		"owner_reviewer": "some (Doc <: owner) & (Doc <: reviewer)",
		"admin_member":   "some (Group <: admin) & (Group <: member)",
		"two_folders":    "some d: Doc | #d.folder > 1",
		"folder_cycle":   "some f: Folder | f in f.^parent",
		"group_cycle":    "some g: Group | g in g.^member_group_member",
	}
	src := "open workspace\nopen workspace_assumptions\n\n" +
		"run satisfiable { assumptions and some Doc.folder and some Folder.parent and some Group.member_group_member } for 6\n"
	for name, violation := range violations {
		src += "run assumed_" + name + " { assumptions and " + violation + " } for 4\n"
		src += "run unassumed_" + name + " { " + violation + " } for 4\n"
	}
	path := filepath.Join(dir, "restrict.als")
	require.NoError(t, os.WriteFile(path, []byte(src), 0644))

	solved := runAlloy(t, path)
	require.True(t, solved["satisfiable"], "assumptions should leave room for instances")
	for name := range violations {
		require.False(t, solved["assumed_"+name], "%v: assumptions should rule it out", name)
		require.True(t, solved["unassumed_"+name], "%v: should be possible without assumptions", name)
	}
}
