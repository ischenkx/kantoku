package jsoncodec

import "encoding/json"

type Dynamic struct {
}

func (d Dynamic) Encode(source any) ([]byte, error) {
	return json.Marshal(source)
}

func (d Dynamic) Decode(payload []byte, destination any) error {
	return json.Unmarshal(payload, destination)
}
