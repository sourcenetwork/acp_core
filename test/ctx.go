package test

import (
	"context"
	"testing"
	"time"

	prototypes "github.com/cosmos/gogoproto/types"
	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/services"
	"github.com/sourcenetwork/acp_core/pkg/types"
	testutil "github.com/sourcenetwork/acp_core/test/util"
	"github.com/stretchr/testify/require"
)

var DefaultTs = testutil.MustDateTimeToProto("2024-01-01 00:00:00")

var _ context.Context = (*TestCtx)(nil)

type TestCtx struct {
	Ctx        context.Context
	T          *testing.T
	Runtime    runtime.RuntimeManager
	Engine     types.ACPEngineServer
	Actors     ActorRegistrar
	State      ActionState
	Playground types.PlaygroundServiceServer
	Time       *prototypes.Timestamp
}

func (t *TestCtx) SetRootPrincipal() {
	t.Ctx = auth.InjectPrincipal(t.Ctx, auth.RootPrincipal())
}

func (t *TestCtx) SetPrincipal(name string) {
	did := t.Actors.DID(name)
	principal, err := auth.NewDIDPrincipal(did)
	require.Nil(t.T, err)
	t.Ctx = auth.InjectPrincipal(t.Ctx, principal)
}

func NewTestCtx(t *testing.T) *TestCtx {
	timeServ := NewConstantTimeService(DefaultTs)
	manager := NewTestRuntime(t, timeServ)
	engine := services.NewACPEngine(manager)

	playground := playgroundFactory(t, manager)

	return &TestCtx{
		Ctx:        context.Background(),
		T:          t,
		Runtime:    manager,
		Engine:     engine,
		Playground: playground,
		Actors: ActorRegistrar{
			actors: make(map[string]string),
		},
		Time: DefaultTs,
	}
}

func (c *TestCtx) Deadline() (deadline time.Time, ok bool) { return c.Ctx.Deadline() }
func (c *TestCtx) Done() <-chan struct{}                   { return c.Ctx.Done() }
func (c *TestCtx) Err() error                              { return c.Ctx.Err() }
func (c *TestCtx) Value(key any) any                       { return c.Ctx.Value(key) }
