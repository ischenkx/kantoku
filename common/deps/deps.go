package deps

import "context"

type Deps interface {
	// return id or struct? struct or *struct?
	Dependency(ctx context.Context, id string) (Dependency, error)
	Resolve(ctx context.Context, id string) error
	Group(ctx context.Context, id string) (Group, error)
	Make(ctx context.Context) (*Dependency, error)                // creates single dependency
	MakeGroup(ctx context.Context, ids ...string) (string, error) // creates group from set of dependencies
	Ready(ctx context.Context) (<-chan string, error)
}
