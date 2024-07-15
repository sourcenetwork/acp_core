package system

import (
	"context"

	rcdb "github.com/sourcenetwork/raccoondb"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

const prefix = "system/"
const key = "params"

type ParamsRepository interface {
	Set(context.Context, *types.Params) error
	GetOrDefault(context.Context) (*types.Params, error)
}

func NewParamsRepository(runtime runtime.RuntimeManager) *KVParamsRepository {
	return &KVParamsRepository{
		runtime: runtime,
	}
}

type KVParamsRepository struct {
	runtime runtime.RuntimeManager
}

func (r *KVParamsRepository) Set(ctx context.Context, params *types.Params) error {
	kv := r.runtime.GetKVStore()
	kv = rcdb.NewWrapperKV(kv, []byte(prefix))

	bytes, err := params.Marshal()
	if err != nil {
		return errors.NewFromCause("could not marshal params:", err, errors.ErrorType_BAD_INPUT)
	}

	err = kv.Set([]byte(key), bytes)
	if err != nil {
		return err
	}
	return nil
}

func (r *KVParamsRepository) GetOrDefault(ctx context.Context) (*types.Params, error) {
	kv := r.runtime.GetKVStore()
	kv = rcdb.NewWrapperKV(kv, []byte(prefix))
	bytes, err := kv.Get([]byte(key))
	if err != nil {
		return nil, err
	}

	if bytes == nil {
		return NewDefaultParams(), nil
	}

	params := &types.Params{}
	err = params.Unmarshal(bytes)
	if err != nil {
		// TODO this might mean data corruption or a bug, should be logged
		return nil, errors.NewFromCause("could not unmarshal stored params, considering setting it", err, errors.ErrorType_INTERNAL)
	}

	return params, nil
}

func NewDefaultParams() *types.Params {
	return &types.Params{}
}
