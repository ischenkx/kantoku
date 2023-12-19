package jsoncodec

import (
	"encoding/json"
	"kantoku/common/codec"
)

var _ codec.Dynamic[[]byte] = (*Dynamic)(nil)

type Dynamic struct {
}

func (d Dynamic) Encode(source any) ([]byte, error) {
	return json.Marshal(source)
}

func (d Dynamic) Decode(payload []byte, destination any) error {
	return json.Unmarshal(payload, destination)
}
