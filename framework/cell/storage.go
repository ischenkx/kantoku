package cell

import (
	"context"
)

type Storage[T any] interface {
	Make(ctx context.Context, data T) (id string, err error)
	Get(ctx context.Context, id string) (Cell[T], error)
}
