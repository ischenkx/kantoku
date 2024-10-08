package manager

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/dependency"
	"github.com/ischenkx/kantoku/pkg/core"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"log/slog"
)

type Manager struct {
	System       core.AbstractSystem
	Dependencies dependency.Manager
	TaskToGroup  TaskToGroup
	Resolvers    map[string]Resolver
	Logger       *slog.Logger
}

type dependenciesInfo struct {
	Specs     []DependencySpec
	GroupID   string
	Instances []string
}

func (manager *Manager) Register(ctx context.Context, id string) error {
	// TODO: use compensating transactions for consistency
	t, err := manager.System.Task(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	var info dependenciesInfo
	if err := mapstructure.Decode(t.Info["dependencies"], &info); err != nil {
		return fmt.Errorf("failed to decode dependencies: %w", err)
	}

	specs := info.Specs

	dependencies, err := manager.Dependencies.NewDependencies(ctx, len(specs))
	if err != nil {
		return fmt.Errorf("failed to allocate new dependencies: %w", err)
	}

	depIDs := lo.Map(dependencies, func(dep dependency.Dependency, _ int) string {
		return dep.ID
	})

	for index, depId := range depIDs {
		spec := specs[index]

		resolver, ok := manager.Resolvers[spec.Name]
		if !ok {
			return fmt.Errorf("failed to find a resolver for '%s'", spec.Name)
		}

		if err := resolver.Bind(ctx, depId, spec.Data); err != nil {
			return fmt.Errorf("failed to bind: %w", err)
		}
	}

	groupId, err := manager.Dependencies.NewGroup(ctx)
	if err != nil {
		return fmt.Errorf("failed to create a dependency group: %w", err)
	}

	if err := manager.TaskToGroup.Save(ctx, id, groupId); err != nil {
		return fmt.Errorf("failed to save a task-group binding: %w", err)
	}

	if err := manager.Dependencies.InitializeGroup(ctx, groupId, depIDs...); err != nil {
		return fmt.Errorf("failed to initialize group: %w", err)
	}

	manager.Logger.Info("saving", slog.Any("instances", depIDs))

	err = manager.System.
		Tasks().
		UpdateByIDs(
			ctx,
			[]string{id},
			map[string]any{
				"info.dependencies.group_id":  groupId,
				"info.dependencies.instances": depIDs,
			},
		)
	if err != nil {
		return fmt.Errorf("failed to update task's info: %w", err)
	}

	return nil
}

func (manager *Manager) Ready(ctx context.Context) (tasks <-chan string, err error) {
	g, ctx := errgroup.WithContext(ctx)

	_tasks := make(chan string, 1024)

	manager.resolveDependencies(ctx, g)

	g.Go(func() error {
		manager.Logger.Debug("collecting ready tasks")
		return manager.collectReadyTasks(ctx, _tasks)
	})

	go func() {
		if err := g.Wait(); err != nil {
			manager.Logger.Info("",
				slog.String("error", err.Error()))
		}
	}()

	return _tasks, nil
}

func (manager *Manager) resolveDependencies(ctx context.Context, g *errgroup.Group) {
	for label, resolver := range manager.Resolvers {
		g.Go(func() error {
			manager.Logger.Info("starting a resolver",
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
					manager.Logger.Debug("received a dependency from resolver",
						slog.String("dependency_id", dep),
						slog.String("resolver", label))
					err := manager.Dependencies.Resolve(ctx, dependency.Dependency{
						ID:     dep,
						Status: dependency.OK,
					})
					if err != nil {
						manager.Logger.Error("failed to resolve a dependency",
							slog.String("error", err.Error()))
					}
				}
			}
			return nil
		})
	}
}

func (manager *Manager) collectReadyTasks(ctx context.Context, tasks chan<- string) error {
	channel, err := manager.Dependencies.ReadyGroups(ctx)
	if err != nil {
		return fmt.Errorf("failed to read ready groups: %w", err)
	}

collector:
	for {
		select {
		case id := <-channel:
			task, err := manager.TaskToGroup.TaskByGroup(ctx, id)
			if err != nil {
				manager.Logger.Error("failed to get task by group",
					slog.String("error", err.Error()))
				continue
			}
			manager.Logger.Debug("received a ready group",
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
