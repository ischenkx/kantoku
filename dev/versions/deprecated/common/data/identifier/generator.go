package identifier

import "context"

type Generator interface {
	New(ctx context.Context) (string, error)
}
