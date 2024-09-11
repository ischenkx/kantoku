package inmem

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ischenkx/kantoku/pkg/common/dependency"
	"github.com/samber/lo"
	"sync"
	"time"
)

var validStatuses = []dependency.Status{
	dependency.OK,
	dependency.Failed,
}

type groupInfo struct {
	deps []string
	done bool
}

type Manager struct {
	dependencies map[string]dependency.Dependency
	groups       map[string]groupInfo
	mu           sync.RWMutex
}

func New() *Manager {
	return &Manager{
		dependencies: make(map[string]dependency.Dependency),
		groups:       make(map[string]groupInfo),
	}
}

func (manager *Manager) LoadDependencies(ctx context.Context, ids ...string) ([]dependency.Dependency, error) {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	ids = lo.Uniq(ids)
	result := make([]dependency.Dependency, 0, len(ids))

	for _, id := range ids {
		if dep, ok := manager.dependencies[id]; ok {
			result = append(result, dep)
		}
	}

	return result, nil
}

func (manager *Manager) LoadGroups(ctx context.Context, ids ...string) ([]dependency.Group, error) {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	ids = lo.Uniq(ids)
	result := make([]dependency.Group, 0, len(ids))

	for _, id := range ids {
		info, ok := manager.groups[id]
		if !ok {
			continue
		}

		group := dependency.Group{
			ID:           id,
			Dependencies: make([]dependency.Dependency, 0, len(info.deps)),
		}
		for _, depId := range info.deps {
			group.Dependencies = append(group.Dependencies, manager.dependencies[depId])
		}
		result = append(result, group)
	}

	return result, nil
}

func (manager *Manager) Resolve(ctx context.Context, values ...dependency.Dependency) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	for _, dep := range values {
		if !lo.Contains(validStatuses, dep.Status) {
			continue
		}
		if _, ok := manager.dependencies[dep.ID]; !ok {
			continue
		}
		manager.dependencies[dep.ID] = dep
	}

	return nil
}

func (manager *Manager) NewDependencies(ctx context.Context, n int) ([]dependency.Dependency, error) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	result := make([]dependency.Dependency, 0, n)
	for i := 0; i < n; i++ {
		dep := dependency.Dependency{
			ID:     uuid.New().String(),
			Status: dependency.Pending,
		}
		manager.dependencies[dep.ID] = dep
		result = append(result, dep)
	}

	return result, nil
}

func (manager *Manager) NewGroup(ctx context.Context, ids ...string) (groupId string, err error) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	for _, id := range ids {
		if _, ok := manager.dependencies[id]; !ok {
			return "", fmt.Errorf("dependency not found: '%s'", id)
		}
	}

	groupId = uuid.New().String()
	manager.groups[groupId] = groupInfo{
		deps: ids,
		done: false,
	}

	return groupId, nil
}

func (manager *Manager) ReadyGroups(ctx context.Context) (<-chan string, error) {
	channel := make(chan string, 1024)

	go manager.pollReadyGroups(ctx, channel)

	go func(ctx context.Context) {
		<-ctx.Done()
		close(channel)
	}(ctx)

	return channel, nil
}

func (manager *Manager) syncGroups() (result []string) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	for id, info := range manager.groups {
		if info.done {
			continue
		}

		resolved := true
		for _, depId := range info.deps {
			if manager.dependencies[depId].Status == dependency.Pending {
				resolved = false
				break
			}
		}

		if !resolved {
			continue
		}

		info.done = true
		manager.groups[id] = info

		result = append(result, id)
	}

	return
}

func (manager *Manager) pollReadyGroups(ctx context.Context, channel chan<- string) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			list := manager.syncGroups()
			for _, id := range list {
				select {
				case <-ctx.Done():
					return
				case channel <- id:
				}
			}
		}
	}
}
