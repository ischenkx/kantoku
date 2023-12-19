package executor

import (
	"context"
)

type Runner interface {
	Run(ctx context.Context, id string) ([]byte, error)
}
