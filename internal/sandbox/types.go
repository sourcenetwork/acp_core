package sandbox

import (
	"fmt"

	rcdb "github.com/sourcenetwork/raccoondb"

	"github.com/sourcenetwork/acp_core/internal/parser"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

// sandboxStorePrefix acts as a namspace for the sandbox repository
const sandboxStorePrefix = "/sandboxes"
const sandboxCounterPrefix = "/sandbox-counter"

// GetManagerForSandbox returns a RuntimeManager which manages the environment of a single sandbox.
// Effectively it wraps the Manager's KVStore by adding the sandbox handle a a prefix,
// this ensures that sandboxes are completely isolated from each other
func GetManagerForSandbox(manager runtime.RuntimeManager, handle uint64) (runtime.RuntimeManager, error) {
	kv := manager.GetKVStore()
	prefix := fmt.Sprintf("/engine/%v", handle)
	kv = rcdb.NewWrapperKV(kv, []byte(prefix))
	engineManager, err := runtime.NewRuntimeManager(
		runtime.WithKVStore(kv),
		runtime.WithLogger(manager.GetLogger()),
	)
	if err != nil {
		return nil, errors.Wrap("manager for sandbox", err)
	}
	return engineManager, nil
}

// parsedSandboxCtx stores all members of SandboxCtx alongside their
// original Locations from the input text buffers.
type parsedSandboxCtx struct {
	PolicyDefinition string
	Relationships    []parser.LocatedObject[*types.Relationship]
	Theorem          *parser.LocatedPolicyTheorem
	Policy           *types.Policy
}

// ToCtx discards Location data and returns a SandboxCtx
func (c *parsedSandboxCtx) ToCtx() *types.SandboxCtx {
	return &types.SandboxCtx{
		Policy:        c.Policy,
		Relationships: utils.MapSlice(c.Relationships, func(o parser.LocatedObject[*types.Relationship]) *types.Relationship { return o.Obj }),
		PolicyTheorem: c.Theorem.ToPolicyTheorem(),
	}
}
