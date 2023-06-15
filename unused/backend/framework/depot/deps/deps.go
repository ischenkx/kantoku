package deps

import (
	"context"
	"kantoku/common/data/transactional"
)

type Deps interface {
	Dependency(ctx context.Context, id string) (Dependency, error)
	Resolve(ctx context.Context, id string) error
	Group(ctx context.Context, id string) (Group, error)
	Make(ctx context.Context) (Dependency, error)                 // creates a single dependency
	MakeGroup(ctx context.Context, ids ...string) (string, error) // creates a group from a set of dependencies
	Ready(ctx context.Context) (<-chan transactional.Object[string], error)
}
