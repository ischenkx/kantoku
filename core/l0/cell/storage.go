package cell

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, data []byte) (string, error)
	Get(ctx context.Context, id string) (Cell, error)
	Set(ctx context.Context, cell Cell) error
	Del(ctx context.Context, id string) error
}
