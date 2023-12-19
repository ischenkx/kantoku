package system

import (
	"github.com/ischenkx/kantoku/pkg/system/kernel/resource"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
)

type TaskInitializer func(task *task.Task)

func WithInputs(resources ...resource.ID) TaskInitializer {
	return func(task *task.Task) {
		task.Inputs = resources
	}
}

func WithOutputs(resources ...resource.ID) TaskInitializer {
	return func(task *task.Task) {
		task.Outputs = resources
	}
}

func WithProperties(properties task.Properties) TaskInitializer {
	return func(task *task.Task) {
		task.Properties = properties
	}
}
