package fn

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core"
	"github.com/ischenkx/kantoku/pkg/core/taskopts"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn/future"
	"github.com/samber/lo"
	"log"
	"reflect"
)

type ScheduledTask struct {
	Type    string
	Inputs  []future.AbstractFuture
	Outputs []future.AbstractFuture
}

type Context struct {
	context.Context
	Scheduled     []ScheduledTask
	FutureStorage future.Storage

	spawnedTasks []string // task ids
}

func NewContext(parent context.Context) *Context {
	return &Context{
		Context:       parent,
		Scheduled:     make([]ScheduledTask, 0),
		FutureStorage: future.NewStorage(),
	}
}

func (ctx *Context) bindObjectToResources(obj any, linkTo []string) ([]future.AbstractFuture, error) {
	arr, err := extractFuturesFromObject(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to extract futures: %w", err)
	}

	if linkTo != nil && len(linkTo) != len(arr) {
		return nil, fmt.Errorf("amount of futures in the object doesn't match amount of resources")
	}

	for i, f := range arr {
		ctx.FutureStorage.AddFuture(f)
		if linkTo != nil {
			res := core.Resource{ID: linkTo[i], Status: core.ResourceStatuses.Allocated}
			ctx.FutureStorage.AssignResource(f, &res, false)
		}
	}
	return arr, nil
}

func (ctx *Context) spawn(sys core.AbstractSystem, parentTask core.Task) error {
	// sort in reverse top-sort order to ensure minimal possible execution while rollback is possible
	for _, t := range ctx.Scheduled {
		fut2res := func(fut future.AbstractFuture, _ int) string {
			return ctx.FutureStorage.GetResource(fut).ID
		}

		inputs := lo.Map(t.Inputs, fut2res)
		outputs := lo.Map(t.Outputs, fut2res)
		deps := lo.Map(inputs, func(res string, _ int) taskopts.Dependency {
			return taskopts.Dependency{
				Name: "resource_db",
				Data: res,
			}
		})

		spawned, err := sys.Spawn(ctx,
			core.New(
				taskopts.WithInputs(inputs...),
				taskopts.WithOutputs(outputs...),
				taskopts.WithProperty("context_parent_id", parentTask.ID),
				taskopts.WithType(t.Type),
				taskopts.WithDependencies(deps...),
				taskopts.WithContextID(parentTask.ContextID()),
			),
		)
		if err != nil {
			return fmt.Errorf("failed to spawn task: %w", err)
		}

		ctx.spawnedTasks = append(ctx.spawnedTasks, spawned.ID)
	}
	return nil
}

func (ctx *Context) rollback(sys core.AbstractSystem, err error) {
	log.Printf("encountered error: %s", err)

	err = ctx.FutureStorage.Rollback(ctx, sys.Resources())
	if err != nil {
		log.Printf("failed to rollback resources: %s", err)
	}

	err = sys.Tasks().Delete(ctx, ctx.spawnedTasks)
	if err != nil {
		log.Printf("failed to rollback spawned tasks: %s", err)
	} else {
		ctx.spawnedTasks = []string{}
	}
}

func extractFuturesFromObject(obj any) ([]future.AbstractFuture, error) {
	value := reflect.ValueOf(obj)
	_type := value.Type()
	arr := make([]future.AbstractFuture, 0, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		if field.Kind() != reflect.Struct {
			return nil, fmt.Errorf("field '%s.%s' is not a future, kind=%s",
				_type.Name(),
				_type.Field(i).Name,
				field.Kind(),
			)
		}

		fut, ok := field.Interface().(future.AbstractFuture)
		if !ok {
			return nil, fmt.Errorf("field '%s.%s' is not a future",
				_type.Name(),
				_type.Field(i).Name,
			)
		}

		arr = append(arr, fut)
	}
	return arr, nil
}
