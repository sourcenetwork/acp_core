package test

import (
	"github.com/sourcenetwork/acp_core/pkg/did"
)

type ActorRegistrar struct {
	actors map[string]string
}

// DID returns gets or creates the DID for a the named Actor
func (r *ActorRegistrar) DID(name string) string {
	didStr, ok := r.actors[name]
	if !ok {
		didStr, _, _ = did.ProduceDID()
		r.actors[name] = didStr
	}
	return didStr
}
