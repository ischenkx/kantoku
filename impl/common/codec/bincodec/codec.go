package bincodec

import (
	"kantoku/common/codec"
)

type Codec struct {
}

var _ codec.Codec[[]byte, []byte] = Codec{}

func (c Codec) Encode(t []byte) ([]byte, error) {
	return t, nil
}

func (c Codec) Decode(data []byte) ([]byte, error) {
	return data, nil
}
