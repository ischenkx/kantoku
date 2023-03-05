package deps

import "context"

type Deps interface {
	Dependency(ctx context.Context, id string) (Dependency, error)
	Resolve(ctx context.Context, id string) error
	Group(ctx context.Context, id string) (Group, error)
	Make(ctx context.Context, ids ...string) (string, error)
	Ready(ctx context.Context) (<-chan string, error)
}
