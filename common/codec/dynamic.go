package codec

type Dynamic[T any] interface {
	DynamicEncoder[T]
	DynamicDecoder[T]
}

type DynamicEncoder[T any] interface {
	Encode(source any) (T, error)
}

type DynamicDecoder[T any] interface {
	Decode(payload T, destination any) error
}
