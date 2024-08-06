package sandbox

import (
	"context"
	"fmt"
	"sync"

	rcdb "github.com/sourcenetwork/raccoondb"

	"github.com/sourcenetwork/acp_core/internal/raccoon"
	"github.com/sourcenetwork/acp_core/pkg/playground"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
)

var _ rcdb.Ider[*playground.SandboxRecord] = (*sandboxIder)(nil)

type sandboxIder struct{}

func (i *sandboxIder) Id(obj *playground.SandboxRecord) []byte {
	return i.HandleToBytes(obj.Handle)
}

func (i *sandboxIder) HandleToBytes(handle uint64) []byte {
	return []byte(fmt.Sprintf("%v", handle))
}

func NewSandboxRepository(manager runtime.RuntimeManager) *SandboxRepository {
	marshaler := raccoon.NewGogoProtoMarshaler(func() *playground.SandboxRecord { return &playground.SandboxRecord{} })
	store := rcdb.NewObjStore(manager.GetKVStore(), marshaler, &sandboxIder{})
	return &SandboxRepository{
		lock:  sync.Mutex{},
		store: &store,
		ider:  sandboxIder{},
	}
}

type SandboxRepository struct {
	// TODO make it thread safe
	lock  sync.Mutex
	store *rcdb.ObjectStore[*playground.SandboxRecord]
	ider  sandboxIder
}

func (r *SandboxRepository) GetSandbox(ctx context.Context, handle uint64) (*playground.SandboxRecord, error) {
	opt, err := r.store.GetObject(r.ider.HandleToBytes(handle))
	if err != nil {
		return nil, err
	}
	if opt.IsEmpty() {
		return nil, nil
	}
	return opt.Value(), nil
}

func (r *SandboxRepository) ListSandboxes(ctx context.Context) ([]*playground.SandboxRecord, error) {
	return r.store.List()
}

func (r *SandboxRepository) SetRecord(ctx context.Context, record *playground.SandboxRecord) error {
	return r.store.SetObject(record)
}
