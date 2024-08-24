package fn

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"reflect"
)

func WithContext[T any](ctx context.Context, sys system.AbstractSystem, f func(*Context) (T, error)) (T, error) {
	var result T
	err := schedulingContext(ctx, sys, func(c *Context) error {
		res, err := f(c)
		if err != nil {
			return err
		}

		result = res
		return nil
	})

	return result, err
}

func schedulingContext(ctx context.Context, sys system.AbstractSystem, f func(*Context) error) error {
	proxy := proxyTask{f: f}
	exe := NewExecutor[proxyTask, EmptyStruct, EmptyStruct](proxy)
	taskCtx := NewContext(ctx)

	sysTask := task.Task{
		Inputs:  []resource.ID{},
		Outputs: []resource.ID{},
		ID:      "proxy---should-not-see-this",
		Info:    map[string]any{},
	}

	input, err := exe.prepareInput(taskCtx, sys, sysTask)
	if err != nil {
		return fmt.Errorf("failed to prepare input: %w", err)
	}

	output, err := proxy.Call(taskCtx, input)
	if err != nil {
		return err
	}

	return exe.save(taskCtx, sys, sysTask, output)
}

type proxyTask struct {
	f func(*Context) error
	Function[proxyTask, EmptyStruct, EmptyStruct]
}

type EmptyStruct struct{}

func (t proxyTask) EmptyOutput() (EmptyStruct, error) {
	return EmptyStruct{}, nil
}

func (t proxyTask) Call(ctx *Context, input EmptyStruct) (EmptyStruct, error) {
	err := t.f(ctx)
	return EmptyStruct{}, err
}

func (t proxyTask) InputType() reflect.Type {
	return reflect.TypeOf(EmptyStruct{})
}
