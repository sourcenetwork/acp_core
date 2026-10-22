package alloy

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	pc "github.com/sourcenetwork/acp_core/internal/policy_compiler"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var update = flag.Bool("update", false, "update golden files")

func Test_Emit_MatchesGoldenFiles(t *testing.T) {
	model, err := pc.Compile(workspacePolicy, types.PolicyMarshalingType_YAML)
	require.NoError(t, err)

	for _, f := range Emit(model, "") {
		path := filepath.Join("testdata", f.Name)
		if *update {
			require.NoError(t, os.MkdirAll("testdata", 0755))
			require.NoError(t, os.WriteFile(path, []byte(f.Content), 0644))
		}
		golden, err := os.ReadFile(path)
		require.NoError(t, err)
		require.Equal(t, string(golden), f.Content, "%v differs from golden file, run with -update to regenerate", f.Name)
	}
}

func Test_ModuleName(t *testing.T) {
	require.Equal(t, "my_policy_v2", moduleName("My Policy-v2"))
	require.Equal(t, "policy_1st", moduleName("1st"))
	require.Equal(t, "policy", moduleName(""))
}
