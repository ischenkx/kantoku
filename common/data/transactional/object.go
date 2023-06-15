package transactional

import "context"

type Object[Item any] interface {
	Get(ctx context.Context) (Item, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
