package taskdep

import (
	"kantoku/kernel"
)

func Dep(task string) kernel.Option {
	return func(ctx *kernel.Context) error {
		data := ctx.Data().GetOrSet("taskdep", func() any { return &PluginData{} }).(*PluginData)
		data.Subtasks = append(data.Subtasks, task)

		return nil
	}
}
