package functional

import (
	"context"
	"errors"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/future"
	"reflect"
)

type Executor[T Task[I, O], I, O any] struct {
	task T
}

func NewExecutor[T Task[I, O], I, O any](t T) Executor[T, I, O] {
	return Executor[T, I, O]{task: t}
}

func (e Executor[T, I, O]) Execute(ctx context.Context, sys system.AbstractSystem, task task.Task) error {
	taskCtx, input, err := e.prepare(ctx, sys, task)

	out, err := e.task.Call(taskCtx, input)
	if err != nil {
		return err
	}

	// here we can check for circular dependencies if we want...

	err = e.save(taskCtx, sys, task, out)

	// write status about successful execution or something

	return err
}

func (e Executor[T, I, O]) Type() string {
	return taskType[I, O](e.task)
}

func (e Executor[T, I, O]) prepare(ctx context.Context, sys system.AbstractSystem, task task.Task) (*Context, I, error) {
	var input I
	inputResources, err := sys.Resources().Load(ctx, task.Inputs...)
	if err != nil {
		return nil, input, err
	}

	taskCtx := NewContext(ctx)

	input, err = e.buildInput(inputResources, taskCtx.FutureStorage)
	if err != nil {
		return nil, input, err
	}
	return taskCtx, input, err
}

func (e Executor[T, I, O]) save(ctx *Context, sys system.AbstractSystem, task task.Task, out O) error {
	ctx.addFutureStruct(out, task.Outputs)
	err := ctx.FutureStorage.Encode(codec.JSON[any]())
	if err != nil {
		return err
	}

	// all futures are created and added to ctx
	err = ctx.FutureStorage.Allocate(ctx, sys.Resources())
	if err != nil {
		ctx.rollback(sys, err)
		return err
	}

	err = ctx.spawn(sys)
	if err != nil {
		ctx.rollback(sys, err)
		return err
	}

	err = ctx.FutureStorage.Save(ctx, sys.Resources())
	if err != nil {
		ctx.rollback(sys, err)
		return err
	}
	return nil
}

// can replace any in return value to 'I', but it's hard to return empty value this way
func (e Executor[T, I, O]) buildInput(resources []resource.Resource, storage future.Storage) (I, error) {
	structType := e.task.InputType()
	structValue := reflect.New(structType).Elem()
	structInterface, ok := structValue.Interface().(I)
	if !ok {
		return structInterface, errors.New("not convertable to input")
	}

	// Get the number of fields in the struct
	numFields := structValue.NumField()
	if numFields != len(resources) {
		return structInterface, errors.New("input struct doesn't match inputs")
	}

	// Initialize struct fields from the fields array
	for i := 0; i < numFields; i++ {
		if resources[i].Status != resource.Ready {
			return structInterface, errors.New("not ready resource at position")
		}

		err := parseField(resources[i].Data, structValue.Field(i))
		if err != nil {
			return structInterface, err
		}

		// save resource to storage so they won't be copied
		fut, ok := structValue.Field(i).Interface().(future.AbstractFuture)
		if !ok {
			return structInterface, errors.New("cannot convert field to future")
		}
		storage.AddFuture(fut)
		storage.AssignResource(fut, &resources[i], true)
	}

	// Return the initialized struct
	structInterface, ok = structValue.Interface().(I)
	if !ok {
		return structInterface, errors.New("not convertable to input")
	}
	return structInterface, nil
}

func parseField(data []byte, field reflect.Value) error {
	uninitializedFut := reflect.New(field.Type())
	futAndErr := uninitializedFut.MethodByName("ParseToNew").Call([]reflect.Value{reflect.ValueOf(data)})
	if !futAndErr[1].IsNil() {
		err, ok := futAndErr[1].Interface().(error)
		if !ok {
			panic("can't convert error to error")
		}
		return err
	}
	if !field.CanSet() {
		return errors.New("cannot set field in input")
	}
	if !futAndErr[0].CanConvert(field.Type()) {
		return errors.New("cannot convert parsed future to field type")
	}
	field.Set(futAndErr[0])
	return nil
}

func taskType[I, O any](task Task[I, O]) string {
	typ := reflect.ValueOf(task).Type()
	return typ.PkgPath() + "." + typ.Name()
}
