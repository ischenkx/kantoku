package bincodec

import "kantoku/common/codec"

var _ codec.Codec[[]byte, []byte] = Codec{}

type Codec struct {
}

func (c Codec) Encode(t []byte) ([]byte, error) {
	return t, nil
}

func (c Codec) Decode(data []byte) ([]byte, error) {
	return data, nil
}
