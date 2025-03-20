package types

import (
	"fmt"

	didpkg "github.com/sourcenetwork/acp_core/pkg/did"
)

// AnonymousPrincipal returns a generic principal representing unauthenticated users
func AnonymousPrincipal() Principal {
	return Principal{
		Kind:       PrincipalKind_Anonymous,
		Identifier: "",
	}
}

// RootPrincipal returns the root for the system
func RootPrincipal() Principal {
	return Principal{
		Kind:       PrincipalKind_Root,
		Identifier: PrincipalKind_Root.String(),
	}
}

// NewDIDPrincipal attempts to create a Principal from a provided did.
// If did is invalid, return an error
func NewDIDPrincipal(did string) (Principal, error) {
	err := didpkg.IsValidDID(did)
	if err != nil {
		return Principal{}, fmt.Errorf("invalid principal id: %v", err)
	}
	return Principal{
		Kind:       PrincipalKind_DID,
		Identifier: did,
	}, nil
}

func (p *Principal) Equals(other *Principal) bool {
	return p.Identifier == other.Identifier && p.Kind == other.Kind
}
