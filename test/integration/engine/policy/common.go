package policy

import "github.com/sourcenetwork/acp_core/pkg/types"

var metadata *types.SuppliedMetadata = &types.SuppliedMetadata{
	Attributes: map[string]string{
		"test": "abc",
	},
	Blob: []byte("test"),
}
