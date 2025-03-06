package ppp

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ Transformer = (*IdTransformer)(nil)

type IdTransformer struct {
	counterVal uint64
}

func NewIdTransformer(policyCounterVal uint64) Transformer {
	return &IdTransformer{
		counterVal: policyCounterVal,
	}
}

func (s *IdTransformer) Validate(pol *types.Policy) *errors.MultiError {
	if pol.Id == "" {
		return errors.NewMultiError(ErrIdTransformer, fmt.Errorf("id not specified"))
	}
	return nil
}

// normalize normalizes a policy by setting default values for optional fields.
func (s *IdTransformer) Transform(producer PolicyProvider) (*types.Policy, *errors.MultiError) {
	pol := producer()

	if pol.ActorResource == nil {
		pol.ActorResource = &types.ActorResource{
			Name: defaultActorResourceName,
		}
	}

	// policy is sorted before building id to ensure determinism
	pol.Sort()
	pol.Id = s.id(pol, s.counterVal)

	return pol, nil
}

// buildId computes the unique id for a policy.
//
// the id is a hash of the policy hash and the policy counter number
func (t *IdTransformer) id(pol *types.Policy, counter uint64) string {
	hasher := sha256.New()

	hasher.Write(t.hashPol(pol))
	hasher.Write([]byte(fmt.Sprintf("%v", counter)))

	hash := hasher.Sum(nil)
	id := hex.EncodeToString(hash)
	return id
}

// hashPol computes a partial sha256 hash of a policy.
// the hashing algorithm includes a subset of the fields which are deterministic.
func (t *IdTransformer) hashPol(pol *types.Policy) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(pol.Name))

	for _, resource := range pol.Resources {
		hasher.Write([]byte(resource.Name))

		for _, rel := range resource.Relations {
			hasher.Write([]byte(rel.Name))
		}

		for _, perm := range resource.Permissions {
			hasher.Write([]byte(perm.Name))
			hasher.Write([]byte(perm.Expression))
		}
	}

	return hasher.Sum(nil)
}
