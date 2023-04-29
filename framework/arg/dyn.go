package arg

import (
	"kantoku"
	"kantoku/impl/common/codec/jsoncodec"
)

func Parse[T any](data []byte) (Dyn[T], error) {
	return jsoncodec.New[Dyn[T]]().Decode(data)
}

type Dyn[T any] struct {
	Data T
}

func (dyn Dyn[T]) Initialize(ctx *kantoku.Context) ([]byte, error) {
	return jsoncodec.New[Dyn[T]]().Encode(dyn)
}

func Arg[T any](arg T) Dyn[T] {
	return Dyn[T]{Data: arg}
}
