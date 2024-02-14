package taskutil

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
		existingDependencies, ok := t.Info["dependencies"].([]map[string]any)
		if !ok {
			existingDependencies = []map[string]any{}
			t.Info["dependencies"] = existingDependencies
		}

		for _, dep := range dependencies {
			existingDependencies = append(existingDependencies, map[string]any{
				"name": dep.Name,
				"data": dep.Data,
			})
		}
		t.Info["dependencies"] = existingDependencies
	}
}
