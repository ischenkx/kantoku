package deps

import (
	"context"
	"kantoku/common/data/transactional"
)

type Deps interface {
	Dependency(ctx context.Context, id string) (Dependency, error)
	Resolve(ctx context.Context, id string) error
	Group(ctx context.Context, id string) (Group, error)
	Make(ctx context.Context) (Dependency, error) // creates a single dependency

	// MakeGroup receives func that is called before group can be resolved (and get to Ready channel)
	// It is preferable to avoid doing changes to Deps state before intercept.
	// Also intercept should not fail anything if MakeGroup fails after its call.
	MakeGroup(ctx context.Context, intercept func(ctx context.Context, id string) error,
		ids ...string) (string, error)
	Ready(ctx context.Context) (<-chan transactional.Object[string], error)
}
