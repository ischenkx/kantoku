package kantoku

import (
	"context"
	"github.com/google/uuid"
	"kantoku/common/data/kv"
	"kantoku/core/l0/event"
	"kantoku/framework/cell"
	"kantoku/framework/depot"
)

type Kantoku struct {
	events event.Bus
	depot  *depot.Depot
	tasks  kv.Database[Task]
	cells  cell.Storage[[]byte]
}

func New(config Config) *Kantoku {
	return &Kantoku{
		events: config.Events,
		depot:  config.Depot,
		tasks:  config.Tasks,
		cells:  config.Cells,
	}
}

func (kantoku *Kantoku) New(ctx context.Context, task Task) (id string, err error) {
	task.Spec.ID = uuid.New().String()

	_, err = kantoku.tasks.Set(ctx, task.Spec.ID, task)
	if err != nil {
		return "", err
	}

	err = kantoku.depot.Schedule(ctx, task.Spec.ID, task.Dependencies)
	if err != nil {
		return "", err
	}

	return task.Spec.ID, nil
}

func (kantoku *Kantoku) Events() event.Bus {
	return kantoku.events
}

func (kantoku *Kantoku) Depot() *depot.Depot {
	return kantoku.depot
}

func (kantoku *Kantoku) Tasks() kv.Reader[Task] {
	return kantoku.tasks
}

func (kantoku *Kantoku) Cells() cell.Storage[[]byte] {
	return kantoku.cells
}
