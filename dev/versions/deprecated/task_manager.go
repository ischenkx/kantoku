package deprecated

import (
	"context"
	"kantoku/framework/job"
	"kantoku/framework/plugins/info"
)

type TaskManager struct {
	kantoku *Kantoku
}

func (manager TaskManager) ByID(id string) *Task {
	return &Task{
		id:      id,
		kantoku: manager.kantoku,
	}
}

func (manager TaskManager) Info() *info.Storage {
	return manager.kantoku.info
}

func (manager TaskManager) Spawn(ctx context.Context, spec Spec) (job.SpawnResult, error) {
	payload, err := manager.kantoku.parametrizationCodec.Encode(spec.parametrization)
	if err != nil {
		return job.SpawnResult{}, err
	}
	var options []job.Option

	options = append(options, withParametrization(spec.parametrization))
	options = append(options, spec.opts...)

	if manager.kantoku.settings.AutoInputDependencies {
		options = append(options, AutoInputDeps())
	}

	return manager.kantoku.jobs.Spawn(ctx, job.Describe(spec.typ, payload).With(options...))
}
