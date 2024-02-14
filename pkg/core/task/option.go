package task

import (
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/core/resource"
)

type Option func(t *Task)

func WithID(id string) Option {
	return func(t *Task) {
		t.ID = id
	}
}

func WithInputs(inputs ...resource.ID) Option {
	return func(t *Task) {
		t.Inputs = inputs
	}
}

func WithOutputs(outputs ...resource.ID) Option {
	return func(t *Task) {
		t.Outputs = outputs
	}
}

func WithProperty(key string, value any) Option {
	return func(t *Task) {
		t.Info[key] = value
	}
}

func WithInfo(info record.R) Option {
	return func(t *Task) {
		if info == nil {
			info = record.R{}
		}
		t.Info = info
	}
}
