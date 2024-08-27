package test

import (
	"context"
	goruntime "runtime"
	"testing"
	"time"

	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/services"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/stretchr/testify/require"
)

var _ context.Context = (*TestCtx)(nil)

type TestCtx struct {
	Ctx        context.Context
	T          *testing.T
	Runtime    runtime.RuntimeManager
	Engine     types.ACPEngineServer
	Actors     ActorRegistrar
	State      ActionState
	Playground types.PlaygroundServiceServer
}

func (t *TestCtx) SetPrincipal(name string) {
	did := t.Actors.DID(name)
	principal, err := auth.NewDIDPrincipal(did)
	require.Nil(t.T, err)
	t.Ctx = auth.InjectPrincipal(t.Ctx, principal)
}

func NewTestCtx(t *testing.T) *TestCtx {
	manager := NewTestRuntime(t)
	engine := services.NewACPEngine(manager)
	playground := services.NewPlaygroundService(manager)

	// if the current runtime is JS, meaning we are running some sort of external executor
	// use the playground JS instead, which wraps the proxy js playground back into
	// a PlaygroundServiceServer implmenetation, while creating objects
	// in the JS context to receive as arguments
	if goruntime.GOOS == "js" {
		t.Log("using JS Playground")
		playground = newPlaygroundJSImpl(t, manager)
	}

	return &TestCtx{
		Ctx:        context.Background(),
		T:          t,
		Runtime:    manager,
		Engine:     engine,
		Playground: playground,
		Actors: ActorRegistrar{
			actors: make(map[string]string),
		},
	}
}

func (c *TestCtx) Deadline() (deadline time.Time, ok bool) { return c.Ctx.Deadline() }
func (c *TestCtx) Done() <-chan struct{}                   { return c.Ctx.Done() }
func (c *TestCtx) Err() error                              { return c.Ctx.Err() }
func (c *TestCtx) Value(key any) any                       { return c.Ctx.Value(key) }
