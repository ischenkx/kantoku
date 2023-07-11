package functions

import "context"

type Runner interface {
	Run(ctx context.Context) error
}
