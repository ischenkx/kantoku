package cell

import (
	"context"
)

// TODO: make implicit id generation for the purpose of immutability
type Storage[T any] interface {
	Make(ctx context.Context, data T) error
	Get(ctx context.Context, id string) (Cell[T], error)
}
