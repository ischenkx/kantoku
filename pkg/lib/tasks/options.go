package tasks

import (
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/samber/lo"
)

type Dependency struct {
	Name string
	Data any
}

func DependOnInputs() task.Option {
	return func(t *task.Task) {
		WithDependencies(
			lo.Map(t.Inputs, func(id string, _ int) Dependency {
				return ResourceDependency(id)
			})...,
		)(t)
	}
}

func ResourceDependency(id string) Dependency {
	return Dependency{
		Name: "resource",
		Data: id,
	}
}

func WithDependencies(dependencies ...Dependency) task.Option {
	return func(t *task.Task) {
		existingDependencies, ok := t.Info["dependencies"].(map[string]any)
		if !ok {
			existingDependencies = map[string]any{}
			t.Info["dependencies"] = existingDependencies
		}

		existingSpecs, ok := existingDependencies["specs"].([]map[string]any)
		if !ok {
			existingSpecs = []map[string]any{}
			existingDependencies["specs"] = existingSpecs
		}

		for _, dep := range dependencies {
			existingSpecs = append(existingSpecs, map[string]any{
				"name": dep.Name,
				"data": dep.Data,
			})
		}
		existingDependencies["specs"] = existingSpecs
		t.Info["dependencies"] = existingDependencies
	}
}

func WithContextID(ctxId string) task.Option {
	return task.WithProperty("context_id", ctxId)
}
