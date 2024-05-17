package future

import (
	"encoding/json"
	"errors"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"sync/atomic"
)

type fid int32

var idCounter fid = 0

// AbstractFuture can hold a future without caring about it's type
type AbstractFuture interface {
	IsFilled() bool
	Encode(codec codec.Codec[any, []byte]) ([]byte, error)
	getId() fid
}

type Future[T any] struct {
	id     fid // local id, valid only during program run
	value  *T
	filled bool
}

func (f Future[T]) getId() fid {
	return f.id
}

func (f Future[T]) Value() T { // or *T?
	return *f.value
}

func (f Future[T]) IsFilled() bool {
	return f.filled
}

func Empty[T any]() Future[T] {
	atomic.AddInt32((*int32)(&idCounter), 1)
	return Future[T]{filled: false, id: fid(atomic.LoadInt32((*int32)(&idCounter)))}
}

func FromValue[T any](val T) Future[T] {
	atomic.AddInt32((*int32)(&idCounter), 1)
	return Future[T]{value: &val, filled: true, id: fid(atomic.LoadInt32((*int32)(&idCounter)))}
}

func (f Future[T]) ParseToNew(data []byte) (Future[T], error) {
	var val T
	err := json.Unmarshal(data, &val)
	if err != nil {
		return Future[T]{}, err
	}

	return FromValue[T](val), nil
}

func (f Future[T]) EmptyValue() T {
	var val T
	return val
}

func (f Future[T]) Encode(codec codec.Codec[any, []byte]) ([]byte, error) {
	if !f.IsFilled() {
		return nil, errors.New("can't make resource from empty future")
	}
	data, err := codec.Encode(f.Value())
	if err != nil {
		return nil, err
	}
	return data, nil
}
