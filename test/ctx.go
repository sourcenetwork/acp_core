package test

import (
	"context"
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
	runtime := NewTestRuntime(t)
	engine := services.NewACPEngine(runtime)
	playground := services.NewPlaygroundService(runtime)
	return &TestCtx{
		Ctx:        context.Background(),
		T:          t,
		Runtime:    runtime,
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
