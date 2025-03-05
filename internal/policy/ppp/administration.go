package ppp

import (
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ Transformer = (*DecentralizedAdminTransformer)(nil)

type DecentralizedAdminTransformer struct{}

func (t *DecentralizedAdminTransformer) Name() string {
	return "Decentralized Administration"
}

func (t *DecentralizedAdminTransformer) Validate(policy *types.Policy) []error {
	var violations []error
	return violations
}

func (t *DecentralizedAdminTransformer) Transform(provider PolicyProvider) (*types.Policy, error) {
	policy := provider()
	return policy, nil
}
