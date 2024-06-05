package access_decision

/*
import (
	"context"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	gogoproto "github.com/cosmos/gogoproto/proto"
	raccoon "github.com/sourcenetwork/raccoondb"

	"github.com/sourcenetwork/acp_core/internal/utils"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func NewGogoProtoMarshaler[T gogoproto.Message](factory func() T) *GogoProtoMarshaler[T] {
	return &GogoProtoMarshaler[T]{
		factory: factory,
	}
}

type GogoProtoMarshaler[T gogoproto.Message] struct {
	factory func() T
}

func (m *GogoProtoMarshaler[T]) Marshal(t *T) ([]byte, error) {
	return gogoproto.Marshal(*t)
}

func (m *GogoProtoMarshaler[T]) Unmarshal(bytes []byte) (T, error) {
	t := m.factory()
	err := gogoproto.Unmarshal(bytes, t)
	if err != nil {
		return t, err
	}

	return t, nil
}

type AccessDecisionRepository struct {
	kv storetypes.KVStore
}

func NewAccessDecisionRepository(store storetypes.KVStore) *AccessDecisionRepository {
	return &AccessDecisionRepository{
		kv: store,
	}
}

func (r *AccessDecisionRepository) getStore(ctx context.Context) raccoon.ObjectStore[*types.AccessDecision] {
	rcKV := raccoon.KvFromCosmosKv(r.kv)
	marshaler := NewGogoProtoMarshaler[*types.AccessDecision](func() *types.AccessDecision { return &types.AccessDecision{} })
	ider := &decisionIder{}
	return raccoon.NewObjStore[*types.AccessDecision](rcKV, marshaler, ider)
}

func (r *AccessDecisionRepository) wrapErr(err error) error {
	if err == nil {
		return err
	}

	return fmt.Errorf("%v: %w", err, types.ErrAcpInternal)
}

func (r *AccessDecisionRepository) Set(ctx context.Context, decision *types.AccessDecision) error {
	store := r.getStore(ctx)
	err := store.SetObject(decision)
	return r.wrapErr(err)
}

func (r *AccessDecisionRepository) Get(ctx context.Context, id string) (*types.AccessDecision, error) {
	store := r.getStore(ctx)
	opt, err := store.GetObject([]byte(id))
	var obj *types.AccessDecision
	if !opt.IsEmpty() {
		obj = opt.Value()
	}
	return obj, r.wrapErr(err)
}

func (r *AccessDecisionRepository) Delete(ctx context.Context, id string) error {
	store := r.getStore(ctx)
	err := store.DeleteById([]byte(id))
	return r.wrapErr(err)
}

func (r *AccessDecisionRepository) ListIds(ctx context.Context) ([]string, error) {
	store := r.getStore(ctx)
	bytesIds, err := store.ListIds()
	ids := utils.MapSlice(bytesIds, func(bytes []byte) string { return string(bytes) })
	return ids, r.wrapErr(err)
}

func (r *AccessDecisionRepository) List(ctx context.Context) ([]*types.AccessDecision, error) {
	store := r.getStore(ctx)
	objs, err := store.List()
	return objs, r.wrapErr(err)
}

type decisionIder struct{}

func (i *decisionIder) Id(decision *types.AccessDecision) []byte {
	return []byte(decision.Id)
}

*/
