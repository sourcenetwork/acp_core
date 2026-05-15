package zanzi

import (
	"encoding/hex"
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestPolicyMapperAppDataMarshalWithPolicyAttributesIsDeterministic(t *testing.T) {
	record := &types.PolicyRecord{
		Policy: &types.Policy{
			Id:            "policy",
			Name:          "Policy",
			ActorResource: &types.ActorResource{Name: "actor"},
			Attributes: map[string]string{
				"alpha":   "one",
				"bravo":   "two",
				"charlie": "three",
				"delta":   "four",
				"echo":    "five",
				"foxtrot": "six",
			},
		},
		Metadata: &types.RecordMetadata{
			Supplied: &types.SuppliedMetadata{
				Attributes: map[string]string{
					"hotel":  "seven",
					"india":  "eight",
					"juliet": "nine",
				},
			},
		},
		PolicyDefinition: "name: Policy",
		MarshalType:      types.PolicyMarshalingType_YAML,
	}

	encodings := map[string]struct{}{}
	mapper := &policyMapper{}
	for i := 0; i < 1000; i++ {
		zanziRecord, err := mapper.ToZanziRecord(record)
		require.NoError(t, err)
		encodings[hex.EncodeToString(zanziRecord.AppData)] = struct{}{}
	}

	require.Len(t, encodings, 1, "policy app data must serialize deterministically")
}
