package system

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	event2 "github.com/ischenkx/kantoku/pkg/system/kernel/event"
	"github.com/ischenkx/kantoku/pkg/system/kernel/resource"
	task2 "github.com/ischenkx/kantoku/pkg/system/kernel/task"
	"github.com/samber/lo"
	"log/slog"
)

var _ AbstractSystem = (*System)(nil)

type System struct {
	events    event2.Bus
	resources resource.Storage
	tasks     task2.Storage
	info      record.Storage
}

func New(
	events event2.Bus,
	resources resource.Storage,
	tasks task2.Storage,
	info record.Storage) *System {

	return &System{
		events:    events,
		resources: resources,
		tasks:     tasks,
		info:      info,
	}
}

func (system *System) Tasks() task2.Storage {
	return system.tasks
}

func (system *System) Resources() resource.Storage {
	return system.resources
}

func (system *System) Events() event2.Bus {
	return system.events
}

// TODO: move all constants (events, "task_id", etc) to actual constant (probably it'd be better
// TODO: to have a separate package for event names, so they can be referred from the kernel and this high-level package

func (system *System) Spawn(ctx context.Context, initializers ...TaskInitializer) (Task, error) {
	var newTask task2.Task
	for _, initializer := range initializers {
		if initializer == nil {
			continue
		}
		initializer(&newTask)
	}

	type state struct {
		task task2.Task
	}

	// TODO: transactions must provide atomicity and eventual consistency guarantees
	// Thus, the transaction below must be executed via Sagas
	// (this pattern is very similar to this code but it allows retries for compensating transactions)
	tx := lo.NewTransaction[state]().
		Then(
			func(s state) (state, error) {
				createdTask, err := system.Tasks().Create(ctx, s.task)
				if err != nil {
					err = fmt.Errorf("failed to create a task: %w", err)
				}
				s.task = createdTask
				return s, err
			},
			func(s state) state {
				if err := system.Tasks().Delete(ctx, s.task.ID); err != nil {
					slog.Error("failed to delete a task in a compensating transaction",
						slog.String("id", s.task.ID),
						slog.String("error", err.Error()))
				}
				return s
			}).
		Then(
			func(state state) (state, error) {
				err := system.Info().Insert(ctx, record.R{
					InfoTaskID: state.task.ID,
				})
				if err != nil {
					return state, err
				}

				return state, nil
			},
			func(state state) state {
				err := system.Info().Filter(record.R{InfoTaskID: state.task.ID}).Erase(ctx)
				if err != nil {
					slog.Error("failed to delete task info in the compensating transaction",
						slog.String("id", state.task.ID),
						slog.String("error", err.Error()))
				}

				return state
			}).
		Then(
			func(s state) (state, error) {
				err := system.Events().Publish(ctx, event2.New(TaskNewEvent, []byte(s.task.ID)))
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
		return Task{}, err
	}

	return Task{ID: result.task.ID, System: system}, nil
}

func (system *System) Task(id string) Task {
	return Task{
		ID:     id,
		System: system,
	}
}

func (system *System) Info() record.Storage {
	return system.info
}
