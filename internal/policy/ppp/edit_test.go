package ppp

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestPreservedResourcesRequirementValidateReturnsSortedErrors(t *testing.T) {
	oldPolicy := &types.Policy{
		ActorResource: &types.ActorResource{Name: "actor"},
		Resources: []*types.Resource{
			{Name: "zulu"},
			{Name: "alpha"},
			{Name: "mike"},
		},
	}
	newPolicy := types.Policy{
		ActorResource: &types.ActorResource{Name: "actor"},
	}

	requirement := NewPreservedResourcesRequirement(oldPolicy)
	for range 100 {
		err := requirement.Validate(newPolicy)
		require.EqualError(t, err,
			"cannot remove resources from Policy; attrs={}; kind=BAD_INPUT: "+
				"removed resource alpha; removed resource mike; removed resource zulu; ")
	}
}
