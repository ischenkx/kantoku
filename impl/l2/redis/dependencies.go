package redis

import (
	"context"
	"kantoku/l2"
)

type Dependencies struct {
}

func (d *Dependencies) Resolve(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (d *Dependencies) Unique(ctx context.Context) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dependencies) Get(ctx context.Context) (l2.Dependency, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dependencies) Resolutions(ctx context.Context) <-chan l2.Dependency {
	//TODO implement me
	panic("implement me")
}
