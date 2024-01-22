package system

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/task"
)

type AbstractSystem interface {
	Tasks() record.Set[task.Task]
	Resources() resource.Storage
	Events() *event.Broker

	Spawn(ctx context.Context, t task.Task) (task.Task, error)
	Task(ctx context.Context, id string) (task.Task, error)
}
