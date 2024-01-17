package system

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/system/kernel/event"
	"github.com/ischenkx/kantoku/pkg/system/kernel/resource"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
)

type AbstractSystem interface {
	Tasks() task.Storage
	Resources() resource.Storage
	Events() event.Bus
	Info() record.Storage
	Spawn(ctx context.Context, initializers ...TaskInitializer) (*Task, error)
	Task(id string) *Task
}
