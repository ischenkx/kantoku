package alpha

import "context"

type Runner interface {
	Run(ctx context.Context, alpha Alpha) ([]byte, error)
}
