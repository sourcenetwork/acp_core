package version

import "github.com/sourcenetwork/acp_core/pkg/types"

var (

	// Version contains the version number for the current acp_core build - if available
	Version string = ""

	// Commit contains the commit hash for the current acp_core build - if available
	Commit string = ""
)

// GetBuildInfo returns the current BuildInfo for this compilation
func GetBuildInfo() types.BuildInfo {
	return types.BuildInfo{
		Version: Version,
		Commit:  Commit,
	}
}
