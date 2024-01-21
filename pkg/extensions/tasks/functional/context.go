package functional

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/extensions/tasks/future"
	"github.com/ischenkx/kantoku/pkg/processors/scheduler/dependencies/simple/manager"
	"github.com/ischenkx/kantoku/pkg/system"
	"github.com/ischenkx/kantoku/pkg/system/kernel/resource"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
	"github.com/samber/lo"
	"log"
	"reflect"
)

type ScheduledTask struct {
	Name    string
	Inputs  []future.Future[any]
	Outputs []future.Future[any]
}

type Context struct {
	context.Context
	Scheduled     []ScheduledTask
	FutureStorage future.Storage

	spawnedLog []string // task ids
}

func NewContext(parent context.Context) Context {
	return Context{Context: parent, Scheduled: make([]ScheduledTask, 0), FutureStorage: future.NewStorage()}
}

func Execute[T Task[I, O], I, O any](ctx Context, task T, input I) O {
	out := task.EmptyOutput()
	ctx.Scheduled = append(ctx.Scheduled, ScheduledTask{
		Name:    taskName[I, O](task),
		Inputs:  ctx.addFutureStruct(input),
		Outputs: ctx.addFutureStruct(out),
	})
	return out
}

// doesn't care about resources!
func (ctx *Context) addFutureStruct(obj any) []future.Future[any] {
	arr := futureStructToArr(obj)
	for _, f := range arr {
		ctx.FutureStorage.AddFuture(f)
	}
	return arr
}

func (ctx *Context) spawn(sys system.AbstractSystem) error {
	for _, t := range ctx.Scheduled {
		fut2res := func(fut future.Future[any], _ int) resource.ID {
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
			system.WithInputs(inputs...),
			system.WithOutputs(outputs...),
			system.WithProperties(
				task.Properties{
					Data: map[string]any{
						"type":         t.Name,
						"dependencies": deps,
					},
				},
			))
		if err != nil {
			return err
		}
		ctx.spawnedLog = append(ctx.spawnedLog, spawned.ID)
		log.Println("Spawned:", spawned)
	}
	return nil
}

func (ctx *Context) rollback(sys system.AbstractSystem) {
	err := ctx.FutureStorage.Rollback(ctx, sys.Resources())
	if err != nil {
		log.Printf("failed to rollback resources: %s", err)
	}

	err = sys.Tasks().Delete(ctx, ctx.spawnedLog...)
	if err != nil {
		log.Printf("failed to rollback spawned tasks: %s", err)
	} else {
		ctx.spawnedLog = []string{}
	}
}

//func (ctx Context) Schedule(storage resource.Storage) error {
//	flatOutputs := lo.FlatMap(ctx.Scheduled, func(item ScheduledTask, _ int) []future.Future[any] {
//		return futureStructToArr(item.Outputs)
//	})
//	uniqOutputs := map[future.Future[any]]any{}
//	for _, o := range flatOutputs {
//		uniqOutputs[o] = nil
//	}
//
//	resIds, err := storage.Alloc(ctx, len(uniqOutputs))
//	if err != nil {
//		return err
//	}
//	zip := lo.Zip2[string, future.Future[any]](resIds, lo.Keys(uniqOutputs))
//	filledZip := lo.Filter(zip, func(item lo.Tuple2[string, future.Future[any]], index int) bool {
//		return item.B.IsFilled()
//	})
//	filledResources :=
//	for _, t := range filledZip {
//		storage.Load(ctx)
//		storage.
//	}
//}

func taskName[I, O any](task Task[I, O]) string {
	return reflect.ValueOf(task).Type().Name()
}

func futureStructToArr(obj any) []future.Future[any] {
	v := reflect.ValueOf(obj)
	arr := make([]future.Future[any], v.NumField())
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fmt.Printf("%s: %v\n", v.Type().Field(i).Name, field.Interface())

		if field.Kind() == reflect.Struct && field.Type().ConvertibleTo(reflect.TypeOf(future.Future[any]{})) {
			x, ok := field.Interface().(future.Future[any])
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

// compare to this:
//func lol() {
//	var resources []resource.Resource
//
//	// Get the number of fields in the struct
//	structValue := reflect.ValueOf(data)
//	numFields := structValue.NumField()
//
//	values := make([][]byte, numFields)
//	// Initialize struct fields from the fields array
//	for i := 0; i < numFields; i++ {
//		// todo: 100 checks for basic types
//		value, err := e.codec.Encode(structValue.Field(i))
//		if err != nil {
//			return err
//		}
//		values[i] = value
//	}
//
//	ids, err := sys.Resources().Alloc(ctx, numFields)
//	if err != nil {
//		return err
//	}
//	resources = make([]resource.Resource, numFields)
//	for i := 0; i < numFields; i++ {
//		resources[i].ID = ids[i]
//		resources[i].Data = values[i]
//	}
//	return sys.Resources().Init(ctx, resources)
//}
