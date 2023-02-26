package l2

import (
	"context"
)

type Scheduler interface {
	Schedule(ctx context.Context, id string) error
	Pending(ctx context.Context) <-chan string
}
