package future

import (
	"errors"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/system/kernel/resource"
	"sync/atomic"
)

type fid int32

var idCounter fid = 0

type Future[T any] struct {
	id     fid // local id, valid only during program run
	value  *T
	filled bool
}

func (f *Future[T]) Value() T { // or *T?
	return *f.value
}

func (f *Future[T]) IsFilled() bool {
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

func FromResource[T any](res resource.Resource, codec codec.Codec[T, []byte]) (Future[T], error) {
	value, err := codec.Decode(res.Data)
	if err != nil {
		return Future[T]{}, err
	}
	return Future[T]{value: &value, filled: true}, nil
}

func (f *Future[T]) ToResource(codec codec.Codec[T, []byte]) (resource.Resource, error) {
	if !f.IsFilled() {
		return resource.Resource{}, errors.New("can't make resource from empty future")
	}
	data, err := codec.Encode(f.Value())
	if err != nil {
		return resource.Resource{}, err
	}
	return resource.Resource{Data: data}, nil
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
