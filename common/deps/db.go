package deps

import "context"

type DB interface {
	Dependency(ctx context.Context, id string) (Dependency, error)
	Resolve(ctx context.Context, id string) error
	Group(ctx context.Context, id string) (Group, error)
	Make(ctx context.Context, deps ...string) (Group, error)
	Ready(ctx context.Context) (<-chan string, error)
}
