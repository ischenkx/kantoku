package taskdep

import (
	job2 "kantoku/framework/job"
)

func Dep(task string) job2.Option {
	return func(ctx *job2.Context) error {
		data := ctx.
			Data().
			GetOrSet("taskdep", func() any { return &PluginData{} }).(*PluginData)
		data.Subtasks = append(data.Subtasks, task)

		return nil
	}
}
