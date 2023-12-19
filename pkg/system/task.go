package system

import (
	"context"
	"errors"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/system/kernel/resource"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
)

const InfoTaskID = "task_id"

type Task struct {
	ID     string
	System AbstractSystem
}

func (t *Task) Info(ctx context.Context) (record.R, error) {
	iter := t.System.
		Info().
		Filter(record.R{InfoTaskID: t.ID}).
		Cursor().
		Iter()
	defer iter.Close(ctx)

	rec, err := iter.Next(ctx)
	if err != nil {
		if errors.Is(err, record.ErrIterEmpty) {
			rec = record.R{}
		} else {
			return nil, err
		}
	}

	return rec, nil
}

func (t *Task) Inputs(ctx context.Context) ([]resource.ID, error) {
	raw, err := t.Raw(ctx)
	if err != nil {
		return nil, err
	}

	return raw.Inputs, nil
}

func (t *Task) Outputs(ctx context.Context) ([]resource.ID, error) {
	raw, err := t.Raw(ctx)
	if err != nil {
		return nil, err
	}

	return raw.Outputs, nil
}

func (t *Task) Properties(ctx context.Context) (task.Properties, error) {
	raw, err := t.Raw(ctx)
	if err != nil {
		return task.Properties{}, err
	}

	return raw.Properties, nil
}

func (t *Task) Raw(ctx context.Context) (task.Task, error) {
	batch, err := t.System.Tasks().Load(ctx, t.ID)
	if err != nil {
		return task.Task{}, err
	}

	return batch[0], nil
}
