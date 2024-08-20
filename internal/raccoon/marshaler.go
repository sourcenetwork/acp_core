// package raccoon contains helpers for raccoondb
package raccoon

import (
	gogoproto "github.com/cosmos/gogoproto/proto"
	raccoon "github.com/sourcenetwork/raccoondb"
)

var _ raccoon.Marshaler[gogoproto.Message] = (*GogoProtoMarshaler[gogoproto.Message])(nil)

// Return a new Marshaler from a factory function which produces container objects of type T
func NewGogoProtoMarshaler[T gogoproto.Message](factory func() T) *GogoProtoMarshaler[T] {
	return &GogoProtoMarshaler[T]{
		factory: factory,
	}
}

// GogoProtoMarshaler implements Raccoon's Marshaler for cosmos gogoproto messages
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
