package raccoon

import (
	"context"
	"encoding/binary"
	"sync"

	rcdb "github.com/sourcenetwork/raccoondb"

	"github.com/sourcenetwork/acp_core/pkg/runtime"
)

// FIXME make this concurrent
const key = "value"

const kvPrefix = "/counter"

// NewCounterStoreFromRuntimeManager builds a CounterStore from resources in the manager
func NewCounterStoreFromRunetimeManager(manager runtime.RuntimeManager, prefix string) CounterStore {
	return NewCounterStore(prefix, manager.GetKVStore(), manager.GetLogger())
}

// NewCounterStore returns a new counter store
func NewCounterStore(prefix string, kv rcdb.KVStore, logger runtime.Logger) CounterStore {
	if prefix != "" {
		kv = rcdb.NewWrapperKV(kv, []byte(prefix))
	}
	kv = rcdb.NewWrapperKV(kv, []byte(kvPrefix))

	return CounterStore{
		kv:     kv,
		logger: logger,
		prefix: prefix,
		lock:   sync.Mutex{},
	}
}

// CounterStore wraps a KV Store and creates a monotomically increasing counter
type CounterStore struct {
	kv     rcdb.KVStore
	logger runtime.Logger
	prefix string
	lock   sync.Mutex
}

func (r *CounterStore) getStore() rcdb.KVStore {
	return r.kv
}

// GetFree returns the next free number in the counter
func (r *CounterStore) GetNext(ctx context.Context) (uint64, error) {
	kv := r.getStore()

	var currID uint64 = 0
	counter, err := kv.Get([]byte(key))
	if err != nil {
		return 0, err
	}
	if counter != nil {
		currID = r.decode(counter)
	}

	freeID := currID + 1
	return freeID, nil
}

// Increment updates the counter to the next free number
func (r *CounterStore) setCounter(ctx context.Context, counter uint64) error {
	kv := r.getStore()

	err := kv.Set([]byte(key), r.encode(counter))
	if err != nil {
		return err
	}

	return nil
}

// encode maps a uin64 to a big endian byte slice
func (r *CounterStore) encode(counter uint64) []byte {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, counter)
	return buff
}

// decode converts a big endian byte slice into a uint64
func (r *CounterStore) decode(counter []byte) uint64 {
	return binary.BigEndian.Uint64(counter)
}

// Increment increments the counter by 1
func (r *CounterStore) Increment(ctx context.Context) error {
	_, err := r.GetNextAndIncrement(ctx)
	return err
}

// GetNextAndIncrement atomically gets the next free counter and increments it
func (r *CounterStore) GetNextAndIncrement(ctx context.Context) (uint64, error) {
	free, err := r.GetNext(ctx)
	if err != nil {
		return 0, err
	}

	err = r.setCounter(ctx, free)
	if err != nil {
		return 0, err
	}

	return free, nil
}

// Acquire does a blocking call to acquire exclusive rights to CounteStore.
// Acquire MUST be followed by a defer call to the returned Releaser `Release`.
// Failing to do so means the store will be indefinitely locked,
// potentially causing a deadlock.
func (r *CounterStore) Acquire() *Releaser {
	r.lock.Lock()
	r.logger.Debug("counter store %v locked", r)
	return &Releaser{
		called: false,
		store:  r,
	}
}

// Releaser models a one-shot callback which can be used to release a store
type Releaser struct {
	store  *CounterStore
	called bool
}

// Release frees up a CounterStore to be used by some other thread.
// Calling release more than once will result in a panic.
func (r *Releaser) Release() {
	if r.called {
		panic("Releaser was previously released!")
	}
	r.store.lock.Unlock()
	r.called = true
	r.store.logger.Debug("counter store %v unlocked", r)
}
