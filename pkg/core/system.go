package core

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/uid"
	"github.com/samber/lo"
	"log/slog"
)

type AbstractSystem interface {
	Tasks() TaskDB
	Resources() ResourceDB
	Events() Broker

	Spawn(ctx context.Context, t Task) (Task, error)
	Task(ctx context.Context, id string) (Task, error)
}

var _ AbstractSystem = (*System)(nil)

type System struct {
	broker    Broker
	resources ResourceDB
	tasks     TaskDB

	logger *slog.Logger
}

func NewSystem(broker Broker, resources ResourceDB, tasks TaskDB, logger *slog.Logger) *System {
	return &System{
		broker:    broker,
		resources: resources,
		tasks:     tasks,
		logger:    logger,
	}
}

func (system System) Tasks() TaskDB {
	return system.tasks
}

func (system System) Resources() ResourceDB {
	return system.resources
}

func (system System) Events() Broker {
	return system.broker
}

func (system System) Spawn(ctx context.Context, newTask Task) (initializedTask Task, err error) {
	type state struct {
		task Task
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
				err := system.Tasks().Insert(ctx, []Task{state.task})
				if err != nil {
					return state, err
				}

				return state, nil
			},
			func(state state) state {
				err := system.tasks.Delete(ctx, []string{state.task.ID})
				if err != nil {
					system.logger.Error("failed to delete task info in the compensating transaction",
						slog.String("id", state.task.ID),
						slog.String("error", err.Error()))
				}

				return state
			}).
		Then(
			func(s state) (state, error) {
				// todo: enable back
				err := system.Events().Send(ctx, NewEvent(OnTask.Created, []byte(s.task.ID)))
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

	return result.task, nil
}

func (system System) Task(ctx context.Context, id string) (Task, error) {
	tasks, err := system.Tasks().ByIDs(ctx, []string{id})
	if err != nil {
		return Task{}, fmt.Errorf("failed to load task: %w", err)
	}

	if len(tasks) == 0 {
		return Task{}, fmt.Errorf("task not found: %s", id)
	}

	return tasks[0], nil
}
