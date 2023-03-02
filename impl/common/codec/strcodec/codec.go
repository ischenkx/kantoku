package strcodec

import "io"

type Codec struct {
}

func (c Codec) Encode(t string) ([]byte, error) {
	return []byte(t), nil
}

func (c Codec) Decode(reader io.Reader) (string, error) {
	raw, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
