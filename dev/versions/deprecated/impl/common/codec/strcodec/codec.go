package strcodec

import "kantoku/common/codec"

var _ codec.Codec[string, []byte] = Codec{}

type Codec struct {
}

func (c Codec) Encode(t string) ([]byte, error) {
	return []byte(t), nil
}

func (c Codec) Decode(data []byte) (string, error) {
	return string(data), nil
}
