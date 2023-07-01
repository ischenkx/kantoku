package deps

import (
	"context"
	"kantoku/common/data/transactional"
)

type Deps interface {
	Dependency(ctx context.Context, id string) (Dependency, error)
	Resolve(ctx context.Context, id string) error
	Group(ctx context.Context, id string) (Group, error)
	MakeDependency(ctx context.Context) (Dependency, error) // creates a single dependency

	// MakeGroupId generates id for a group, which then can be passed to SaveGroup
	MakeGroupId(ctx context.Context) (string, error)
	// SaveGroup saves group with given id and dependencies to Deps
	SaveGroup(ctx context.Context, groupId string, depIds ...string) error
	Ready(ctx context.Context) (<-chan transactional.Object[string], error)
}
