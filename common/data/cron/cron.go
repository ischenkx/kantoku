package cron

import (
	"context"
	"time"
)

type Cron interface {
	Schedule(ctx context.Context, at time.Time, event string) error
	Events(ctx context.Context) (<-chan string, error)
}
