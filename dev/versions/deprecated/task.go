package deprecated

import (
	"context"
	"fmt"
	"kantoku/common/data"
	"kantoku/common/data/deps"
	"kantoku/common/data/future"
	"kantoku/framework/job"
	"kantoku/framework/plugins/info"
	"kantoku/framework/plugins/status"
)

type Task struct {
	id              string
	kantoku         *Kantoku
	parametrization Parametrization
	raw             job.Job
	cachedRaw       bool
	cachedParam     bool
}

func (task *Task) ID() string {
	return task.id
}

func (task *Task) Status(ctx context.Context) (status.Status, error) {
	value, err := task.Info().Get(ctx, "status")
	if err != nil {
		if err == data.NotFoundErr {
			return status.Unknown, nil
		}
		return "", fmt.Errorf("failed to retrieve status: %s", err)
	}

	stat, ok := value.(status.Status)
	if !ok {
		return status.Unknown, fmt.Errorf("failed to cast retrieved value to a status struct (value='%s')", value)
	}

	return stat, nil
}

func (task *Task) Context(ctx context.Context) (string, error) {
	value, err := task.Info().Get(ctx, "context")
	if err != nil {
		if err == data.NotFoundErr {
			return "", nil
		}
		return "", fmt.Errorf("failed to retrieve status: %s", err)
	}

	cont, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("failed to cast retrieved value to a context struct (value='%s')", value)
	}

	return cont, nil
}

func (task *Task) Info() info.Info {
	return task.kantoku.Info().Get(task.id)
}

func (task *Task) Dependencies(ctx context.Context) ([]deps.Dependency, error) {
	groupID, err := task.kantoku.Dependencies().Depot().GroupTaskBimap().ByValue(ctx, task.ID())
	if err != nil {
		return nil, err
	}

	group, err := task.kantoku.Dependencies().Depot().Deps().Group(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return group.Dependencies, err
}

func (task *Task) Static(ctx context.Context) ([]byte, error) {
	err := task.loadParametrization(ctx)
	if err != nil {
		return nil, err
	}

	return task.parametrization.Static, nil
}

func (task *Task) Inputs(ctx context.Context) ([]future.ID, error) {
	err := task.loadParametrization(ctx)
	if err != nil {
		return nil, err
	}

	return task.parametrization.Inputs, nil
}

func (task *Task) Outputs(ctx context.Context) ([]future.ID, error) {
	err := task.loadParametrization(ctx)
	if err != nil {
		return nil, err
	}

	return task.parametrization.Outputs, nil
}

func (task *Task) Type(ctx context.Context) (string, error) {
	err := task.loadRaw(ctx)
	if err != nil {
		return "", err
	}

	return task.raw.Type, nil
}

func (task *Task) Kantoku() *Kantoku {
	return task.kantoku
}

func (task *Task) loadParametrization(ctx context.Context) error {
	if task.cachedParam {
		return nil
	}
	if err := task.loadRaw(ctx); err != nil {
		return err
	}

	param, err := task.kantoku.parametrizationCodec.Decode(task.raw.Data)
	if err != nil {
		return err
	}

	task.cachedParam = true
	task.parametrization = param
	return nil
}

func (task *Task) loadRaw(ctx context.Context) error {
	if task.cachedRaw {
		return nil
	}
	raw, err := task.kantoku.jobs.Task(task.id).Instance(ctx)
	if err != nil {
		return err
	}
	task.raw = raw
	return nil
}
