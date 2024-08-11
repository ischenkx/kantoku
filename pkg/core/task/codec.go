package task

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type Codec struct{}

func (c Codec) Encode(task Task) (map[string]any, error) {
	return task.AsDoc(), nil
}

func (c Codec) Decode(rec map[string]any) (Task, error) {
	if rec == nil {
		return Task{}, errors.New("record is nil")
	}
	var t Task

	if err := mapstructure.Decode(rec, &t); err != nil {
		return Task{}, fmt.Errorf("failed to decode: %w", err)
	}

	return t, nil
}
