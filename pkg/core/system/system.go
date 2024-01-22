package system

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	recutil "github.com/ischenkx/kantoku/pkg/common/data/record/util"
	"github.com/ischenkx/kantoku/pkg/common/data/uid"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/system/events"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/samber/lo"
	"log/slog"
)

// TODO: UPDATE record.Storage, so It supports auto-encoding/decoding of record types

var _ AbstractSystem = (*System)(nil)

type System struct {
	events    *event.Broker
	resources resource.Storage
	tasks     record.Storage[task.Task]
}

func New(events *event.Broker, resources resource.Storage, tasks record.Storage[task.Task]) *System {
	return &System{
		events:    events,
		resources: resources,
		tasks:     tasks,
	}
}

func (system *System) Tasks() record.Set[task.Task] {
	return system.tasks
}

func (system *System) Resources() resource.Storage {
	return system.resources
}

func (system *System) Events() *event.Broker {
	return system.events
}

// TODO: move all constants (events, "task_id", etc) to actual constant (probably it'd be better
// TODO: to have a separate package for event names, so they can be referred from the kernel and this high-level package

func (system *System) Spawn(ctx context.Context, newTask task.Task) (initializedTask task.Task, err error) {
	type state struct {
		task task.Task
	}

	newTask.ID = uid.Generate()

	// TODO: transactions must provide atomicity and eventual consistency guarantees
	// Thus, the transaction below must be executed via Sagas
	// (this pattern is very similar to this code but it allows retries for compensating transactions)
	tx := lo.NewTransaction[state]().
		Then(
			func(state state) (state, error) {
				err := system.tasks.Insert(ctx, newTask)
				if err != nil {
					return state, err
				}

				return state, nil
			},
			func(state state) state {
				err := system.tasks.Filter(record.R{"id": state.task.ID}).Erase(ctx)
				if err != nil {
					slog.Error("failed to delete task info in the compensating transaction",
						slog.String("id", state.task.ID),
						slog.String("error", err.Error()))
				}

				return state
			}).
		Then(
			func(s state) (state, error) {
				err := system.Events().Send(ctx, event.New(events.OnTask.Created, []byte(s.task.ID)))
				if err != nil {
					return s, fmt.Errorf("failed to publish an event: %w", err)
				}
				return s, nil
			},
			func(s state) state {
				return s
			})

	result, err := tx.Process(state{task: newTask})
	if err != nil {
		return task.Task{}, err
	}

	return result.task, nil
}

func (system *System) Task(ctx context.Context, id string) (task.Task, error) {
	t, err := recutil.Single(
		ctx,
		system.
			Tasks().
			Filter(record.R{"id": id}).
			Cursor().
			Iter(),
	)
	if err != nil {
		return task.Task{}, fmt.Errorf("failed to load task: %w", err)
	}

	return t, nil
}
