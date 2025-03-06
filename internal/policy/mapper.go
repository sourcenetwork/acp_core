package policy

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/policy/ppp"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func mapIRIntoPolicy(ir PolicyIR, counter uint64) (*types.Policy, error) {
	policy := &types.Policy{
		Id:            "",
		Name:          ir.Name,
		Description:   ir.Description,
		Attributes:    ir.Attributes,
		Resources:     ir.Resources,
		ActorResource: ir.ActorResource,
	}

	transformers := []ppp.Transformer{
		&ppp.BasicTransformer{},

		&ppp.DiscretionaryTransformer{},

		&ppp.DecentralizedAdminTransformer{},

		ppp.NewIdTransformer(counter),
	}

	specs := []ppp.Specification{}

	pipeline := ppp.NewPipeline(specs, transformers)
	policy, err := pipeline.Process(policy)
	if err != nil {
		return nil, fmt.Errorf("CreatePolicy: %w", err)
	}
	return policy, nil
}
