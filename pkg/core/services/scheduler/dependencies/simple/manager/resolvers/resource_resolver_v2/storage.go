package resourceResolverV2

import "context"

type BindingStorage interface {
	Save(ctx context.Context, dependencyId string, resourceId string) error
	Load(ctx context.Context, resourceId string) (dependencyId string, err error)
}
