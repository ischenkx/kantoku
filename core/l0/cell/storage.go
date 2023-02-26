package cell

import (
	"context"
	"io"
)

type Storage interface {
	Create(ctx context.Context, data []byte) (string, error)
	Get(ctx context.Context, id string) (Cell, error)
	Set(ctx context.Context, cell Cell) error
	Delete(ctx context.Context, id string) error

	io.Closer
}
