package dependency

import (
	"context"
)

// TODO: Add options to NewDependency
// TODO: Add methods "Fail" and "OK"
// TODO: Make groups Fail-able

type Manager interface {
	LoadDependencies(ctx context.Context, ids ...string) ([]Dependency, error)
	LoadGroups(ctx context.Context, ids ...string) ([]Group, error)
	Resolve(ctx context.Context, values ...Dependency) error
	NewDependencies(ctx context.Context, n int) ([]Dependency, error)
	// NewGroup generates id for a group, which then can be passed to SaveGroup
	NewGroup(ctx context.Context) (groupId string, err error)
	InitializeGroup(ctx context.Context, groupId string, dependencyIds ...string) error
	ReadyGroups(ctx context.Context) (<-chan string, error)
}
