package codec

type Codec[From, To any] interface {
	Encode(From) (To, error)
	Decode(To) (From, error)
}
