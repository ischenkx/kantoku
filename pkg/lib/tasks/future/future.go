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
	//Value() any
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
	return Future[T]{filled: false, id: idCounter}
}

func FromValue[T any](val T) Future[T] {
	atomic.AddInt32((*int32)(&idCounter), 1)
	return Future[T]{value: &val, filled: true, id: idCounter}
}

func (f Future[T]) ParseToNew(data []byte) (Future[T], error) {
	var val T
	err := json.Unmarshal(data, &val)
	if err != nil {
		return Future[T]{}, err
	}
	atomic.AddInt32((*int32)(&idCounter), 1)
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

// banned because new Futures should be created. Otherwise, you can easily fill outputs 💀
// func (f *Future[T]) Fill(val T) {
//	if f.filled {
//		panic("tried to set value second time")
//	}
//	f.value = &val
//	f.filled = true
//}

// do we need it outside?
//func (f *Future[T]) Id() FutureId { // or *T?
//	return f.fid
//}

// feels very wrong
//func (f *Future[T]) SetResource(res resource.Resource, codec codec.Codec[T, []byte]) error {
//	if f.HasResource() {
//		panic("setting resource twice")
//	}
//	if res.Data != nil && f.filled {
//		panic("setting filled resource on filled future")
//	}
//	f.resourceId = res.ID
//	if f.filled {
//		data, err := codec.Encode(*f.value)
//		if err != nil {
//			return err
//		}
//		res.Data = data
//	}
//	return nil
//}
