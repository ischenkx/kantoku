package pool

import (
	"context"
	"kantoku/common/data/transactional"
)

type Reader[Item any] interface {
	Read(ctx context.Context) (<-chan transactional.Object[Item], error)
}

// Writer probably should have NewTransaction method
type Writer[Item any] interface {
	// Write *must* write all items in a transaction!
	Write(ctx context.Context, items ...Item) error
}

type Pool[Item any] interface {
	Reader[Item]
	Writer[Item]
}
