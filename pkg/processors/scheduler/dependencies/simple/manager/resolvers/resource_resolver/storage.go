package resourceResolver

import "context"

type Binding struct {
	DependencyId string
	ResourceId   string
}

type Storage interface {
	Save(ctx context.Context, dependencyId string, resourceId string) error
	Resolve(ctx context.Context, resourceIds ...string) error
	Poll(ctx context.Context, limit int) ([]Binding, error)
}
