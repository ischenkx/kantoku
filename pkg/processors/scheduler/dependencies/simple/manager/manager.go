package manager

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/deps"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/system"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"log/slog"
)

type Manager struct {
	system       system.AbstractSystem
	dependencies deps.Manager
	task2group   TaskToGroup
	resolvers    map[string]Resolver
}

func New(
	system system.AbstractSystem,
	dependencies deps.Manager,
	task2group TaskToGroup,
	resolvers map[string]Resolver) *Manager {

	return &Manager{
		system:       system,
		dependencies: dependencies,
		task2group:   task2group,
		resolvers:    resolvers,
	}
}

func (manager *Manager) Register(ctx context.Context, id string) error {
	// TODO: use compensating transactions for consistency
	task, err := manager.system.Task(id).Raw(ctx)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	rawSpecs := task.Properties.Data["dependencies"]

	var specs []DependencySpec
	if err := mapstructure.Decode(rawSpecs, &specs); err != nil {
		return fmt.Errorf("failed to decode dependencies: %w", err)
	}

	dependencies, err := manager.dependencies.NewDependencies(ctx, len(specs))
	if err != nil {
		return fmt.Errorf("failed to allocate new dependencies: %w", err)
	}

	depIds := lo.Map(dependencies, func(dep deps.Dependency, _ int) string {
		return dep.ID
	})

	for index, depId := range depIds {
		spec := specs[index]

		resolver, ok := manager.resolvers[spec.Name]
		if !ok {
			return fmt.Errorf("failed to find a resolver for '%s'", spec.Name)
		}

		if err := resolver.Bind(ctx, depId, spec.Data); err != nil {
			return fmt.Errorf("failed to bind: %w", err)
		}
	}

	groupId, err := manager.dependencies.NewGroup(ctx)
	if err != nil {
		return fmt.Errorf("failed to create a dependency group: %w", err)
	}

	if err := manager.task2group.Save(ctx, id, groupId); err != nil {
		return fmt.Errorf("failed to save a task-group binding: %w", err)
	}

	if err := manager.dependencies.InitializeGroup(ctx, groupId, depIds...); err != nil {
		return fmt.Errorf("failed to initialize group: %w", err)
	}

	err = manager.system.
		Info().
		Filter(record.R{system.InfoTaskID: id}).
		Update(ctx, record.R{"dependencies": depIds, "group_id": groupId}, record.R{system.InfoTaskID: id})
	if err != nil {
		return fmt.Errorf("failed to update task's info: %w", err)
	}

	slog.Info("registered task",
		slog.String("id", id),
		slog.String("group_id", groupId),
		slog.Any("deps", depIds))

	return nil
}

func (manager *Manager) Ready(ctx context.Context) (tasks <-chan string, err error) {
	g, ctx := errgroup.WithContext(ctx)

	_tasks := make(chan string, 1024)

	manager.resolveDependencies(ctx, g)

	g.Go(func() error {
		slog.Info("collecting ready tasks")
		return manager.collectReadyTasks(ctx, _tasks)
	})

	go func() {
		if err := g.Wait(); err != nil {
			slog.Info("failure",
				slog.String("error", err.Error()))
		}
	}()

	return _tasks, nil
}

func (manager *Manager) resolveDependencies(ctx context.Context, g *errgroup.Group) {
	for label, resolver := range manager.resolvers {
		g.Go(func() error {
			slog.Info("starting a resolver",
				slog.String("label", label))

			depsChannel, err := resolver.Ready(ctx)
			if err != nil {
				return err
			}

		depResolver:
			for {
				select {
				case <-ctx.Done():
					break depResolver
				case dep := <-depsChannel:
					slog.Info("received a dependency from resolver",
						slog.String("dependency_id", dep),
						slog.String("resolver", label))
					err := manager.dependencies.Resolve(ctx, deps.Dependency{
						ID:     dep,
						Status: deps.OK,
					})
					if err != nil {
						slog.Info("failed to resolve a dependency",
							slog.String("error", err.Error()))
					}
				}
			}
			return nil
		})
	}
}

func (manager *Manager) collectReadyTasks(ctx context.Context, tasks chan<- string) error {
	channel, err := manager.dependencies.ReadyGroups(ctx)
	if err != nil {
		return fmt.Errorf("failed to read ready groups: %w", err)
	}

collector:
	for {
		select {
		case id := <-channel:
			task, err := manager.task2group.TaskByGroup(ctx, id)
			if err != nil {
				slog.Info("failed to get task by group",
					slog.String("error", err.Error()))
				continue
			}
			slog.Info("received a ready group",
				slog.String("group_id", id),
				slog.String("task_id", task))

			select {
			case <-ctx.Done():
				break collector
			case tasks <- task:
			}
		case <-ctx.Done():
			break collector
		}
	}

	return nil
}
