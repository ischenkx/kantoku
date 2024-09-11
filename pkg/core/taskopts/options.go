package taskopts

import (
	"github.com/ischenkx/kantoku/pkg/core"
	"github.com/samber/lo"
)

type Dependency struct {
	Name string
	Data any
}

func WithInputs(inputs ...string) core.Option {
	return func(t *core.Task) {
		t.Inputs = inputs
	}
}

func WithOutputs(outputs ...string) core.Option {
	return func(t *core.Task) {
		t.Outputs = outputs
	}
}

func WithProperty(key string, value any) core.Option {
	return func(t *core.Task) {
		t.Info[key] = value
	}
}

func WithInfo(info map[string]any) core.Option {
	return func(t *core.Task) {
		if info == nil {
			info = make(map[string]any)
		}
		t.Info = info
	}
}

func DependOnInputs() core.Option {
	return func(t *core.Task) {
		WithDependencies(
			lo.Map(t.Inputs, func(id string, _ int) Dependency {
				return ResourceDependency(id)
			})...,
		)(t)
	}
}

func ResourceDependency(id string) Dependency {
	return Dependency{
		Name: "resource_db",
		Data: id,
	}
}

func WithDependencies(dependencies ...Dependency) core.Option {
	return func(t *core.Task) {
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

func WithContextID(ctxId string) core.Option {
	return WithProperty("context_id", ctxId)
}

func WithType(t string) core.Option {
	return WithProperty("type", t)
}
