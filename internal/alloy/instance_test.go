package alloy

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	pc "github.com/sourcenetwork/acp_core/internal/policy_compiler"
	"github.com/sourcenetwork/acp_core/internal/sandbox"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func Test_ParseInstance_TranslatesToRelationships(t *testing.T) {
	model, err := pc.Compile(workspacePolicy, types.PolicyMarshalingType_YAML)
	require.NoError(t, err)
	data, err := os.ReadFile("testdata/workspace_instance.xml")
	require.NoError(t, err)

	inst, err := ParseInstance(model, data)

	require.NoError(t, err)
	require.Equal(t, "Run nested_viewer for 3 but 2 Actor", inst.Command)
	require.Equal(t, []string{
		"folder:folder0#owner@did:user:alice",
		"group:group0#owner@did:user:alice",
		"folder:folder0#editor@did:user:alice",
		"folder:folder0#parent@folder:folder0",
		"folder:folder0#viewer@group:group0#member",
		"group:group0#admin@did:user:alice",
		"group:group0#member@group:group0#member",
	}, inst.Relationships)
	require.Equal(t, map[string][]string{"$nested_viewer_f": {"folder:folder0"}}, inst.Skolems)
	require.Equal(t, []string{"folder:folder0", "group:group0"}, inst.Objects)
	require.Equal(t, []string{"did:user:alice"}, inst.Actors)
}

func Test_ParseInstance_RejectsInstanceFromAnotherPolicy(t *testing.T) {
	model, err := pc.Compile(sandbox.Samples[0].Data.PolicyDefinition, types.PolicyMarshalingType_YAML)
	require.NoError(t, err)
	data, err := os.ReadFile("testdata/workspace_instance.xml")
	require.NoError(t, err)

	_, err = ParseInstance(model, data)

	require.ErrorContains(t, err, "was the instance generated from this policy?")
}
