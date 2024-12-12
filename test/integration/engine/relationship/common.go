package relationship

import "github.com/sourcenetwork/acp_core/pkg/types"

var metadata *types.SuppliedMetadata = &types.SuppliedMetadata{
	Attributes: map[string]string{
		"test": "abc",
	},
	Misc: []byte("test"),
}
