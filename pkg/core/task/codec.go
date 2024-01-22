package task

import (
	"errors"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/mitchellh/mapstructure"
)

type Codec struct{}

func (c Codec) Encode(task Task) (record.R, error) {
	return record.R{
		"id":      task.ID,
		"inputs":  task.Inputs,
		"outputs": task.Outputs,
		"info":    task.Info,
	}, nil
}

func (c Codec) Decode(rec record.R) (Task, error) {
	if rec == nil {
		return Task{}, errors.New("record is nil")
	}
	var t Task

	if err := mapstructure.Decode(rec, &t); err != nil {
		return Task{}, fmt.Errorf("failed to decode: %w", err)
	}

	return t, nil
}
