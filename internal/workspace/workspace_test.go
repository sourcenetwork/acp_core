package workspace

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

func write(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
}

func Test_Load(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, PolicyFile, "name: notes\n")
	write(t, dir, ConstraintsFile, "assumptions:\n- at_most_one: note.folder\n")
	write(t, dir, ContractsFile, `
requirements:
- name: notes
  statement: there are notes
  contracts:
  - name: notes-exist
    description: Notes can exist.
    properties:
    - type: pred
      description: some note exists
      formula: |
        some Note
`)
	write(t, dir, "worlds/b.yaml", "relationships: |\n  note:n1#owner@did:user:bob\n")
	write(t, dir, "worlds/a.yaml", "relationships: |\n  note:n1#owner@did:user:alice\nassertions: |\n  Authorizations {}\n")
	write(t, dir, "worlds/notes.txt", "ignored")

	ws, err := Load(dir)
	require.NoError(t, err)
	require.Equal(t, "name: notes\n", ws.SystemModel.Policy)
	require.Len(t, ws.SystemModel.Constraints, 1)
	require.Len(t, ws.SystemModel.Requirements, 1)
	require.Equal(t, types.PropertyType_PRED, ws.SystemModel.Requirements[0].Contracts[0].Properties[0].Type)
	require.Len(t, ws.Worlds, 2)
	require.Equal(t, "a", ws.Worlds[0].Name)
	require.Equal(t, "Authorizations {}\n", ws.Worlds[0].Assertions)
	require.Equal(t, "note:n1#owner@did:user:bob\n", ws.Worlds[1].Relationships)
}

func Test_Load_OnlyPolicyIsRequired(t *testing.T) {
	dir := t.TempDir()
	_, err := Load(dir)
	require.Error(t, err)

	write(t, dir, PolicyFile, "name: notes\n")
	ws, err := Load(dir)
	require.NoError(t, err)
	require.Empty(t, ws.SystemModel.Constraints)
	require.Empty(t, ws.SystemModel.Requirements)
	require.Empty(t, ws.Worlds)
}

func Test_Load_InvalidFilesNameTheFile(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, PolicyFile, "name: notes\n")
	write(t, dir, "worlds/bad.yaml", "relationship: oops\n")
	_, err := Load(dir)
	require.ErrorContains(t, err, "worlds/bad.yaml")
}

func Test_MarshalWorld_RoundTrips(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, PolicyFile, "name: notes\n")
	data, err := MarshalWorld(&types.World{Name: "x", Relationships: "note:n1#owner@did:user:bob\nnote:n1#reader@did:user:carol\n"})
	require.NoError(t, err)
	write(t, dir, "worlds/x.yaml", string(data))

	ws, err := Load(dir)
	require.NoError(t, err)
	require.Equal(t, "note:n1#owner@did:user:bob\nnote:n1#reader@did:user:carol\n", ws.Worlds[0].Relationships)
}

func Test_ReadFiles(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, PolicyFile, "name: notes\n")
	write(t, dir, "worlds/a.yaml", "not: [valid")

	files, err := ReadFiles(dir)
	require.NoError(t, err)
	require.Equal(t, "name: notes\n", files.Policy)
	require.Empty(t, files.Contracts)
	require.Equal(t, []WorldFile{{Name: "a", Content: "not: [valid"}}, files.Worlds, "files are read, not parsed")

	write(t, dir, ContractsFile, "requirements: []\n")
	changed, err := ReadFiles(dir)
	require.NoError(t, err)
	require.NotEqual(t, files.Version, changed.Version)

	same, err := ReadFiles(dir)
	require.NoError(t, err)
	require.Equal(t, changed.Version, same.Version)
}
