package jsoncodec

import (
	"bytes"
	"encoding/json"
	"kantoku/common/codec"
)

type Codec[T any] struct{}

var _ codec.Codec[int, []byte] = Codec[int]{}

func New[T any]() Codec[T] {
	return Codec[T]{}
}

func (c Codec[T]) Encode(value T) ([]byte, error) {
	return json.Marshal(value)
}

func (c Codec[T]) Decode(reader []byte) (T, error) {
	var value T
	err := json.NewDecoder(bytes.NewReader(reader)).Decode(&value)
	return value, err
}
