package zanzi

/*

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/internal/auth_engine"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func setup(t *testing.T) (context.Context, auth_engine.AuthEngine) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())
	kv := stateStore.GetCommitKVStore(storeKey)
	engine, err := NewZanzi(kv, log.NewNopLogger())
	require.Nil(t, err)

	return context.Background(), engine
}

func TestZanzi(t *testing.T) {
	auth_engine.TestAuthEngineImpl(t, setup)
}

*/
