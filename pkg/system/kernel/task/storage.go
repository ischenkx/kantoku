package task

import "context"

type Storage interface {
	Create(ctx context.Context, task Task) (Task, error)
	Delete(ctx context.Context, ids ...string) error
	Load(ctx context.Context, ids ...string) ([]Task, error)
}
