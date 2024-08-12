package system

import (
	"context"
	"fmt"
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
	Events_    *event.Broker
	Resources_ resource.Storage
	Tasks_     task.Storage
	Logger     *slog.Logger
}

func (system System) Tasks() task.Storage {
	return system.Tasks_
}

func (system System) Resources() resource.Storage {
	return system.Resources_
}

func (system System) Events() *event.Broker {
	return system.Events_
}

func (system System) Spawn(ctx context.Context, newTask task.Task) (initializedTask task.Task, err error) {
	type state struct {
		task task.Task
	}

	// shallow copying the info to avoid modification of the original object
	shallowCopiedInfo := make(map[string]any)
	for key, val := range newTask.Info {
		shallowCopiedInfo[key] = val
	}
	newTask.Info = shallowCopiedInfo

	// initializing the execution context
	if _, ok := newTask.Info["context_id"]; !ok {
		newTask.Info["context_id"] = uid.Generate()
	}

	newTask.ID = uid.Generate()

	// TODO: transactions must provide atomicity and eventual consistency guarantees
	// Thus, the transaction below must be executed via Sagas
	// (this pattern is very similar to this code but it allows retries for compensating transactions)
	tx := lo.NewTransaction[state]().
		Then(
			func(state state) (state, error) {
				err := system.Tasks().Insert(ctx, []task.Task{state.task})
				if err != nil {
					return state, err
				}

				return state, nil
			},
			func(state state) state {
				err := system.Tasks_.Delete(ctx, []string{state.task.ID})
				if err != nil {
					system.Logger.Error("failed to delete task info in the compensating transaction",
						slog.String("id", state.task.ID),
						slog.String("error", err.Error()))
				}

				return state
			}).
		Then(
			func(s state) (state, error) {
				// todo: enable back
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

func (system System) Task(ctx context.Context, id string) (task.Task, error) {
	tasks, err := system.Tasks().ByIDs(ctx, []string{id})
	if err != nil {
		return task.Task{}, fmt.Errorf("failed to load task: %w", err)
	}

	if len(tasks) == 0 {
		return task.Task{}, fmt.Errorf("task not found: %s", id)
	}

	return tasks[0], nil
}
