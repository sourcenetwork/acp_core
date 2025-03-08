package ppp

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/transformer"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ transformer.Transformer = (*IdTransformer)(nil)

// ErIdTransformer is the root error for IdTransformer
var ErrIdTransformer = errors.New("id policy transformer", errors.ErrorType_BAD_INPUT)

// IdTransformer transforms a Policy by generating its ID
// The ID transformation is done deterministically by hashing a subset
// of the Policy's fields and adding a sequence number
// to ensure equal policies get different Ids.
type IdTransformer struct {
	counterVal uint64
}

// NewIdTransformer returns a new instance of IdTransformer
func NewIdTransformer(policyCounterVal uint64) transformer.Transformer {
	return &IdTransformer{
		counterVal: policyCounterVal,
	}
}

// Validate asserts that the Policy Id has been generated
func (t *IdTransformer) Validate(pol types.Policy) *errors.MultiError {
	if pol.Id == "" {
		return errors.NewMultiError(ErrIdTransformer, fmt.Errorf("id not specified"))
	}
	return nil
}

// Transform produces the policy Id and sets it in the struct
func (t *IdTransformer) Transform(pol types.Policy) (types.Policy, error) {
	hasher := sha256.New()

	hasher.Write(t.hashPol(&pol))
	hasher.Write([]byte(fmt.Sprintf("%v", t.counterVal)))

	hash := hasher.Sum(nil)
	id := hex.EncodeToString(hash)
	pol.Id = id

	return pol, nil
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

func (t *IdTransformer) GetBaseError() error { return ErrIdTransformer }
