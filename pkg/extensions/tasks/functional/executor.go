package functional

import (
	"context"
	"errors"
	codec "github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/extensions/tasks/future"
	"github.com/ischenkx/kantoku/pkg/system"
	"github.com/ischenkx/kantoku/pkg/system/kernel/resource"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
	"reflect"
)

type Executor struct {
	Task  Task[any, any]
	codec codec.Codec[any, []byte]
}

func (e Executor) Execute(ctx context.Context, sys system.AbstractSystem, task task.Task) error {
	inputResources, err := sys.Resources().Load(ctx, task.Inputs...)
	if err != nil {
		return err
	}

	taskCtx := NewContext(ctx)

	input, err := e.buildInput(inputResources, taskCtx.FutureStorage)
	if err != nil {
		return err
	}

	out, err := e.Task.Call(taskCtx, input)
	if err != nil {
		return err
	}

	// here we can check for circular dependencies if we want...

	taskCtx.addFutureStruct(out)
	// all futures are created and added to taskCtx
	err = taskCtx.FutureStorage.Allocate(ctx, sys.Resources())
	if err != nil {
		taskCtx.rollback(sys)
		return err
	}
	err = taskCtx.FutureStorage.Save(ctx, sys.Resources())
	if err != nil {
		taskCtx.rollback(sys)
		return err
	}

	err = taskCtx.spawn(sys)
	if err != nil {
		taskCtx.rollback(sys)
		return err
	}
	return nil
}

func (e Executor) buildInput(resources []resource.Resource, storage future.Storage) (any, error) {
	structType := e.Task.InputType()
	structValue := reflect.New(structType).Elem()

	// Get the number of fields in the struct
	numFields := structValue.NumField()
	if numFields != len(resources) {
		return nil, errors.New("input struct doesn't match inputs")
	}

	// Initialize struct fields from the fields array
	for i := 0; i < numFields; i++ {
		// convert to fut
		fut, err := future.FromResource(resources[i], e.codec)
		if err != nil {
			return nil, err
		}

		// save resource to storage so they won't be copied
		storage.AddFuture(fut)
		storage.AssignResource(fut, &resources[i], true)
		// Set the value of the struct field
		structValue.Field(i).Set(reflect.ValueOf(fut))
	}

	// Return the initialized struct
	return structValue, nil
}
