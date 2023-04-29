package codec

type Codec[From, To any] interface {
	Encoder[From, To]
	Decoder[From, To]
}

type Encoder[From, To any] interface {
	Encode(From) (To, error)
}

type Decoder[From, To any] interface {
	Decode(reader To) (From, error)
}
