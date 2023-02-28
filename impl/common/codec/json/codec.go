package json

import (
	"encoding/json"
	"io"
)

type Codec[T any] struct{}

func New[T any]() Codec[T] {
	return Codec[T]{}
}

func (c Codec[T]) Encode(value T) ([]byte, error) {
	return json.Marshal(value)
}

func (c Codec[T]) Decode(reader io.Reader) (T, error) {
	var value T
	return value, json.NewDecoder(reader).Decode(&value)
}
