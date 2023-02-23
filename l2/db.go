package l2

import "context"

type DB[T Task] interface {
	Get(ctx context.Context, id string) (T, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, task T) (T, error)
	Schedule(ctx context.Context, id string) error
	Pending(ctx context.Context) <-chan T
}
