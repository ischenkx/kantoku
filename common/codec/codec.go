package codec

import "io"

type Codec[T any] interface {
	Encoder[T]
	Decoder[T]
}

type Encoder[T any] interface {
	Encode(T) ([]byte, error)
}

type Decoder[T any] interface {
	Decode(reader io.Reader) (T, error)
}
