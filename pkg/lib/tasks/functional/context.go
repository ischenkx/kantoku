package functional

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/common/data/record/ops"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/services/scheduler/dependencies/simple/manager"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/future"
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

	spawnedLog []string // task ids
}

func NewContext(parent context.Context) *Context {
	return &Context{
		Context:       parent,
		Scheduled:     make([]ScheduledTask, 0),
		FutureStorage: future.NewStorage(),
	}
}

// where do you get task?! - we can remove it and create empty one with reflect
func Execute[T Task[I, O], I, O any](ctx *Context, task T, input I) O {
	out := task.EmptyOutput()
	ctx.Scheduled = append(ctx.Scheduled, ScheduledTask{
		Type:    taskType[I, O](task),
		Inputs:  ctx.addFutureStruct(input, nil),
		Outputs: ctx.addFutureStruct(out, nil),
	})
	return out
}

// doesn't care about resources!
func (ctx *Context) addFutureStruct(obj any, linkTo []resource.ID) []future.AbstractFuture {
	arr := futureStructToArr(obj)
	for i, f := range arr {
		ctx.FutureStorage.AddFuture(f)
		if linkTo != nil {
			res := resource.Resource{ID: linkTo[i], Status: resource.Allocated}
			ctx.FutureStorage.AssignResource(f, &res, false)
		}
	}
	return arr
}

func (ctx *Context) spawn(sys system.AbstractSystem) error {
	// sort in reverse top-sort order to ensure minimal possible execution while rollback is possible
	for _, t := range ctx.Scheduled {
		fut2res := func(fut future.AbstractFuture, _ int) resource.ID {
			return ctx.FutureStorage.GetResource(fut).ID
		}

		inputs := lo.Map(t.Inputs, fut2res)
		outputs := lo.Map(t.Outputs, fut2res)
		deps := lo.Map(inputs, func(res resource.ID, _ int) manager.DependencySpec {
			return manager.DependencySpec{
				Name: "resource",
				Data: res,
			}
		})

		spawned, err := sys.Spawn(ctx,
			task.Task{
				Inputs:  inputs,
				Outputs: outputs,
				Info: record.R{
					"type":         t.Type,
					"dependencies": deps,
				},
			})
		if err != nil {
			return err
		}
		ctx.spawnedLog = append(ctx.spawnedLog, spawned.ID)
	}
	return nil
}

func (ctx *Context) rollback(sys system.AbstractSystem, err error) {
	log.Printf("encountered error: %s", err)

	err = ctx.FutureStorage.Rollback(ctx, sys.Resources())
	if err != nil {
		log.Printf("failed to rollback resources: %s", err)
	}

	err = sys.Tasks().Filter(record.R{"id": ops.In(ctx.spawnedLog)}).Erase(ctx)
	if err != nil {
		log.Printf("failed to rollback spawned tasks: %s", err)
	} else {
		ctx.spawnedLog = []string{}
	}
}

func futureStructToArr(obj any) []future.AbstractFuture {
	v := reflect.ValueOf(obj)
	arr := make([]future.AbstractFuture, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		//fmt.Printf("%s: %v\n", v.Type().Field(i).Type, field.Interface())

		if field.Kind() == reflect.Struct {
			x, ok := field.Interface().(future.AbstractFuture)
			if !ok {
				panic("your struct is still shit")
			}
			arr[i] = x
		} else {
			panic("your struct is shit")
		}
	}
	return arr
}
