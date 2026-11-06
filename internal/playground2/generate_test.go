package playground2

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	pc "github.com/sourcenetwork/acp_core/internal/policy_compiler"
	"github.com/sourcenetwork/acp_core/internal/simulator"
	"github.com/sourcenetwork/acp_core/internal/telemetry"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

const folderPolicy = `
name: folders
resources:
- name: folder
  relations:
  - name: parent
    types: [folder]
  - name: viewer
    types: [actor]
  permissions:
  - name: view
    expr: viewer + parent->view
- name: doc
  relations:
  - name: folder
    types: [folder]
  - name: viewer
    types: [actor]
  permissions:
  - name: view
    expr: viewer + folder->view
`

const folderConstraints = `
assumptions:
- acyclic: folder.parent
- at_most_one: doc.folder
- disjoint: [doc.owner, doc.viewer]
`

func analyzer(t *testing.T) *Analyzer {
	t.Helper()
	if _, err := exec.LookPath("alloy"); err != nil {
		t.Skip("alloy not found in PATH")
	}
	return &Analyzer{AlloyPath: "alloy", Timeout: time.Minute}
}

func folderModel(t *testing.T) *types.SystemModel {
	t.Helper()
	constraints, err := pc.UnmarshalAssumptions(folderConstraints)
	require.NoError(t, err)
	return &types.SystemModel{Policy: folderPolicy, Constraints: constraints}
}

func Test_Generate_ConfigDefaultsAndLimits(t *testing.T) {
	cfg, err := withDefaults(nil)
	require.NoError(t, err)
	require.Equal(t, uint32(defaultCount), cfg.Count)
	require.Equal(t, uint32(defaultScope), cfg.Scope)
	require.Equal(t, uint32(defaultActors), cfg.Actors)

	_, err = withDefaults(&types.GenerateConfig{Count: maxCount + 1})
	require.ErrorIs(t, err, errors.ErrorType_BAD_INPUT)
	_, err = withDefaults(&types.GenerateConfig{Scope: maxScope + 1})
	require.ErrorIs(t, err, errors.ErrorType_BAD_INPUT)
	_, err = withDefaults(&types.GenerateConfig{Actors: maxActors + 1})
	require.ErrorIs(t, err, errors.ErrorType_BAD_INPUT)

	g := &Analyzer{AlloyPath: "alloy", Timeout: time.Minute}
	_, err = g.Generate(context.Background(), &types.GenerateRequest{})
	require.ErrorIs(t, err, ErrGenerate)
}

func Test_Generate_ReturnsValidWorlds(t *testing.T) {
	g := analyzer(t)
	req := &types.GenerateRequest{
		SystemModel: folderModel(t),
		Config:      &types.GenerateConfig{Count: 3, Actors: 4},
	}

	resp, err := g.Generate(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, resp.Worlds, 3)

	manager, err := runtime.NewRuntimeManager()
	require.NoError(t, err)
	for _, world := range resp.Worlds {
		// every actor takes part in some relationship
		actors := map[string]bool{}
		for _, line := range strings.Split(strings.TrimSpace(world.Relationships), "\n") {
			if _, subject, _ := strings.Cut(line, "@"); strings.HasPrefix(subject, "did:") {
				actors[subject] = true
			}
		}
		require.Len(t, actors, 4, world.Relationships)

		// constraints hold: no doc in two folders, no doc owner also listed as a viewer
		folders := map[string]int{}
		owners := map[string]string{}
		for _, line := range strings.Split(strings.TrimSpace(world.Relationships), "\n") {
			object, rest, _ := strings.Cut(line, "#")
			relation, subject, _ := strings.Cut(rest, "@")
			switch {
			case strings.HasPrefix(object, "doc:") && relation == "folder":
				folders[object]++
				require.LessOrEqual(t, folders[object], 1, world.Relationships)
			case strings.HasPrefix(object, "doc:") && relation == "owner":
				owners[object] = subject
			case strings.HasPrefix(object, "doc:") && relation == "viewer":
				require.NotEqual(t, owners[object], subject, world.Relationships)
			}
		}

		// the engine accepts the world
		eval, err := simulator.HandleEvaluate(telemetry.InitRequestContext(context.Background()), manager, &types.EvaluateRequest{
			Policy:        folderPolicy,
			Relationships: world.Relationships,
		})
		require.NoError(t, err)
		require.True(t, eval.Valid, "%v\n%v", eval.Errors, world.Relationships)
	}

	again, err := g.Generate(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, resp, again, "generate should be deterministic")
}
