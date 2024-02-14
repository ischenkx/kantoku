package resource

import "context"

type Storage interface {
	Load(ctx context.Context, ids ...ID) ([]Resource, error)
	Alloc(ctx context.Context, amount int) ([]ID, error)
	Init(ctx context.Context, resources []Resource) error
	Dealloc(ctx context.Context, ids []ID) error
}
