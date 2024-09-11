package fn

import (
	"context"
	"errors"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/core"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn/future"
	"reflect"
	"strings"
)

type Executor[T AbstractFunction[I, O], I, O any] struct {
	task T
}

func NewExecutor[T AbstractFunction[I, O], I, O any](t T) Executor[T, I, O] {
	return Executor[T, I, O]{task: t}
}

func (e Executor[T, I, O]) Execute(ctx context.Context, sys core.AbstractSystem, task core.Task) error {
	taskCtx := NewContext(ctx)

	input, err := e.prepareInput(taskCtx, sys, task)
	if err != nil {
		return fmt.Errorf("failed to prepare input: %w", err)
	}

	out, err := e.task.Call(taskCtx, input)
	if err != nil {
		// TODO: should the err be wrapped
		return err
	}

	// here we can check for circular dependencies if we want...

	err = e.save(taskCtx, sys, task, out)

	// write status about successful execution or something

	return err
}

func (e Executor[T, I, O]) Type() string {
	return taskType[T, I, O]()
}

func (e Executor[T, I, O]) prepareInput(ctx *Context, sys core.AbstractSystem, task core.Task) (I, error) {
	var input I

	inputResources, err := sys.Resources().Load(ctx, task.Inputs...)
	if err != nil {
		return input, fmt.Errorf("failed to load resources: %w", err)
	}

	input, err = e.buildInput(ctx, inputResources)
	if err != nil {
		return input, err
	}

	return input, nil
}

func (e Executor[T, I, O]) save(ctx *Context, sys core.AbstractSystem, task core.Task, out O) error {
	if _, err := ctx.bindObjectToResources(out, task.Outputs); err != nil {
		return fmt.Errorf("failed to bind outputs: %w", err)
	}

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

	err = ctx.spawn(sys, task)
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
func (e Executor[T, I, O]) buildInput(ctx *Context, resources []core.Resource) (I, error) {
	// TODO: use (var input I; reflect.TypeOf(input)
	structType := e.task.InputType()
	structValue := reflect.New(structType).Elem()
	input, ok := structValue.Interface().(I)
	if !ok {
		return input, errors.New("not convertable to input")
	}

	// Get the number of fields in the struct
	numFields := structValue.NumField()
	if numFields != len(resources) {
		return input, errors.New("input struct doesn't match inputs")
	}

	// Initialize struct fields from the fields array
	for i := 0; i < numFields; i++ {
		if resources[i].Status != core.ResourceStatuses.Ready {
			return input, fmt.Errorf("not ready resource_db at position %d", i)
		}

		err := parseField(resources[i].Data, structValue.Field(i))
		if err != nil {
			return input, err
		}

		// save resource_db to storage so they won't be copied
		fut, ok := structValue.Field(i).Interface().(future.AbstractFuture)
		if !ok {
			return input, errors.New("cannot convert field to future")
		}
		ctx.FutureStorage.AddFuture(fut)
		ctx.FutureStorage.AssignResource(fut, &resources[i], true)
	}

	// Return the initialized struct
	input, ok = structValue.Interface().(I)
	if !ok {
		return input, errors.New("not convertable to input")
	}
	return input, nil
}

func parseField(data []byte, field reflect.Value) error {
	uninitializedFut := reflect.New(field.Type())
	futAndErr := uninitializedFut.MethodByName("ParseToNew").Call([]reflect.Value{reflect.ValueOf(data)})
	if !futAndErr[1].IsNil() {
		err, ok := futAndErr[1].Interface().(error)
		if !ok {
			return fmt.Errorf("failed to assert error")
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

func taskType[T AbstractFunction[I, O], I, O any]() string {
	var t T
	typ := reflect.ValueOf(t).Type()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	result := typ.PkgPath() + "/" + typ.Name()

	// TODO: remove:

	result = strings.TrimPrefix(result, "github.com/ischenkx/kantoku/cmd/stand/")

	return result
}
