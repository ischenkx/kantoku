package manager

import "context"

type BindingResult struct {
	Disabled bool
}

type Resolver interface {
	Bind(ctx context.Context, id string, data any) (BindingResult, error)
	Ready(ctx context.Context) (<-chan string, error)
}
