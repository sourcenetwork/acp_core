package raccoon

import (
	"context"
	"strconv"

	rcdb "github.com/sourcenetwork/raccoondb"

	"github.com/sourcenetwork/acp_core/pkg/runtime"
)

// FIXME make this concurrent
const key = "value"

const prefixAugment = "/counter"

func NewCounterStore(manager runtime.RuntimeManager, prefix string) CounterStore {
	return CounterStore{
		runtime: manager,
		prefix:  prefix + prefixAugment,
	}
}

type CounterStore struct {
	runtime runtime.RuntimeManager
	prefix  string
}

func (r *CounterStore) getStore() rcdb.KVStore {
	kv := r.runtime.GetKVStore()
	kv = rcdb.NewWrapperKV(kv, []byte(r.prefix))
	return kv
}

// GetFree returns the next free number in the counter
func (r *CounterStore) GetNext(ctx context.Context) (uint64, error) {
	kv := r.getStore()

	var currID uint64 = 1
	counter, err := kv.Get([]byte(key))
	if err != nil {
		return 0, err
	}
	if counter != nil {
		counterStr := string(counter)
		currID, err = strconv.ParseUint(counterStr, 10, 64)
		if err != nil {
			return 0, err
		}
	}

	freeID := currID + 1
	return freeID, nil
}

// Increment updates the counter to the next free number
func (r *CounterStore) setCounter(ctx context.Context, counter uint64) error {
	kv := r.getStore()

	err := kv.Set([]byte(key), []byte(strconv.FormatUint(counter, 10)))
	if err != nil {
		return err
	}

	return nil
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
