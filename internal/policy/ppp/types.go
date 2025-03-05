package ppp

import "github.com/sourcenetwork/acp_core/pkg/types"

type PolicyProvider func() *types.Policy

type Specification interface {
	Name() string
	Validate(policy *types.Policy) []error
}

type Transformer interface {
	// extends specifier to guarantee the applied transformation wasn't undone / corrupted by some other transformer
	Specification
	Transform(provider PolicyProvider) (*types.Policy, error)
}

// Defra Pol Spec
// specifier
// all resources contains permissions read / write

// Discretionary transformer
// all resources contain owner relation
// computer userset owner is top level node of all trees

// Decentralized admin transform
// forall resource.permission there exists a _can_manage_{} permission
// which includees the onwer

// pipeline: tranform -> apply specification -> model check -> verify

// TODO some field in the policy DSL for a conrtact or something.
// need to specify the policy is used in defra. maybe append only semantics to specifcations / transformers.
