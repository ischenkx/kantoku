package resource

import "context"

type Storage interface {
	Load(ctx context.Context, ids ...string) ([]Resource, error)
	Alloc(ctx context.Context, amount int) ([]string, error)
	Init(ctx context.Context, resources []Resource) error
	Dealloc(ctx context.Context, ids []string) error
}
