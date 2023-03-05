package bincodec

import "io"

type Codec struct {
}

func (c Codec) Encode(t []byte) ([]byte, error) {
	return t, nil
}

func (c Codec) Decode(reader io.Reader) ([]byte, error) {
	return io.ReadAll(reader)
}
