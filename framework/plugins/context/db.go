package context

import "context"

type Database interface {
	Get(ctx context.Context, id string) (Context, error)
	Make(ctx context.Context, parent string) (Context, error)
}
