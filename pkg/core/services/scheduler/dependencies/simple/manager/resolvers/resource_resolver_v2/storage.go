package resourceResolverV2

import "context"

type Binding struct {
	DependencyId string
	ResourceId   string
}

type Storage interface {
	Save(ctx context.Context, dependencyId string, resourceId string) error
	PrepareForResolution(ctx context.Context, resourceIds ...string) error
	Resolve(ctx context.Context, resourceIds ...string) error
}
