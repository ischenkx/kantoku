package executor

import (
	"context"
	"kantoku/common/data/pool"
	job2 "kantoku/framework/job"
)

type Executor struct {
	kernel  *job2.Manager
	runner  Runner
	plugins []Plugin
}

func New(kernel *job2.Manager, runner Runner) *Executor {
	return &Executor{
		kernel: kernel,
		runner: runner,
	}
}

func (executor *Executor) Use(plugins ...Plugin) *Executor {
	executor.plugins = append(executor.plugins, plugins...)
	return executor
}

func (executor *Executor) Run(ctx context.Context) error {
	return pool.AutoCommit[string](ctx, executor.kernel.Inputs(), func(ctx context.Context, id string) error {
		for _, plugin := range executor.plugins {
			if receivedTaskPlugin, ok := plugin.(ReceivedTaskPlugin); ok {
				receivedTaskPlugin.ReceivedTask(ctx, id)
			}
		}

		output, err := executor.runner.Run(ctx, id)
		result := job2.Result{TaskID: id}

		if err != nil {
			result.Data = []byte(err.Error())
			result.Status = job2.FAILURE
		} else {
			result.Data = output
			result.Status = job2.OK
		}

		for _, plugin := range executor.plugins {
			if executedTaskPlugin, ok := plugin.(ExecutedTaskPlugin); ok {
				executedTaskPlugin.ExecutedTask(ctx, id, result)
			}
		}

		err = executor.kernel.Outputs().Set(ctx, result.TaskID, result)
		if err != nil {
			for _, plugin := range executor.plugins {
				if failedToSaveTaskResultPlugin, ok := plugin.(FailedToSaveTaskResultPlugin); ok {
					failedToSaveTaskResultPlugin.FailedToSaveTaskResult(ctx, id, result)
				}
			}
		} else {
			for _, plugin := range executor.plugins {
				if savedTaskResultPlugin, ok := plugin.(SavedTaskResultPlugin); ok {
					savedTaskResultPlugin.SavedTaskResult(ctx, id, result)
				}
			}
		}

		return nil
	})
}
