package strcodec

type Codec struct {
}

func (c Codec) Encode(t string) ([]byte, error) {
	return []byte(t), nil
}

func (c Codec) Decode(data []byte) (string, error) {
	return string(data), nil
}
