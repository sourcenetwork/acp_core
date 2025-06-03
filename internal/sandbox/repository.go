package sandbox

import (
	"context"
	"encoding/binary"

	rcdb "github.com/sourcenetwork/raccoondb"

	"github.com/sourcenetwork/acp_core/internal/raccoon"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var _ rcdb.Ider[*types.SandboxRecord] = (*sandboxIder)(nil)

// sandboxIder implements Raccoon's Ider interface for a SandboxRecord
// Basically serializes a uint64 into a [8]byte
type sandboxIder struct{}

func (i *sandboxIder) Id(obj *types.SandboxRecord) []byte {
	return i.HandleToBytes(obj.Handle)
}

func (i *sandboxIder) HandleToBytes(handle uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, handle)
	return bytes
}

// NewSandboxRepository returns a new Sandbox from a KVStore
func NewSandboxRepository(kv rcdb.KVStore) *SandboxRepository {
	marshaler := raccoon.NewGogoProtoMarshaler(func() *types.SandboxRecord { return &types.SandboxRecord{} })
	store := rcdb.NewObjStore(kv, marshaler, &sandboxIder{})
	return &SandboxRepository{
		store: &store,
		ider:  sandboxIder{},
	}
}

// SandboxRepository manages the creation and retrieval of Sandbox instances
type SandboxRepository struct {
	store *rcdb.ObjectStore[*types.SandboxRecord]
	ider  sandboxIder
}

func (r *SandboxRepository) GetSandbox(ctx context.Context, handle uint64) (*types.SandboxRecord, error) {
	opt, err := r.store.GetObject(r.ider.HandleToBytes(handle))
	if err != nil {
		return nil, errors.NewWithCause("reading sandbox", err, errors.ErrorType_INTERNAL,
			errors.Pair("handle", itoa(handle)),
		)
	}
	if opt.IsEmpty() {
		return nil, nil
	}
	return opt.Value(), nil
}

func (r *SandboxRepository) ListSandboxes(ctx context.Context) ([]*types.SandboxRecord, error) {
	records, err := r.store.List()
	if err != nil {
		return nil, errors.NewWithCause("error loading sandboxes", err, errors.ErrorType_INTERNAL)
	}

	return records, nil
}

func (r *SandboxRepository) SetRecord(ctx context.Context, record *types.SandboxRecord) error {
	err := r.store.SetObject(record)
	if err != nil {
		return errors.NewWithCause("error storing sandbox record", err, errors.ErrorType_INTERNAL,
			errors.Pair("handle", itoa(record.Handle)),
		)
	}
	return nil
}
