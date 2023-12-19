package codec

import "encoding/json"

type Json[From any] struct{}

func (Json[From]) Encode(value From) ([]byte, error) {
	return json.Marshal(value)
}

func (Json[From]) Decode(encoded []byte) (From, error) {
	var value From
	err := json.Unmarshal(encoded, &value)
	return value, err
}

func JSON[From any]() Json[From] {
	return Json[From]{}
}
