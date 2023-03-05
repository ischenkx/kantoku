package task

import "context"

type AbstractTask interface {
	ID(ctx context.Context) string
	Type(ctx context.Context) string
	Argument(ctx context.Context) []byte
}
