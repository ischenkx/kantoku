package manager

import "context"

type Resolver interface {
	Bind(ctx context.Context, id string, data any) error
	Ready(ctx context.Context) (<-chan string, error)
}
