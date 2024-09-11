package core

import "context"

type Resource struct {
	Data   []byte
	ID     string
	Status string
}

type ResourceDB interface {
	Load(ctx context.Context, ids ...string) ([]Resource, error)
	Alloc(ctx context.Context, amount int) (ids []string, err error)
	Init(ctx context.Context, resources []Resource) error
	Dealloc(ctx context.Context, ids []string) error
}
