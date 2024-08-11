package fn_d

import (
	"context"
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

	sysTask := task.Task{
		Inputs:  []resource.ID{},
		Outputs: []resource.ID{},
		ID:      "proxy---should-not-see-this",
		Info:    map[string]any{},
	}

	taskCtx, inp, err := exe.prepare(ctx, sys, sysTask)
	if err != nil {
		return err
	}
	out, err := proxy.Call(taskCtx, inp)
	if err != nil {
		return err
	}

	return exe.save(taskCtx, sys, sysTask, out)
}

type proxyTask struct {
	f func(*Context) error
	Function[proxyTask, EmptyStruct, EmptyStruct]
}

type EmptyStruct struct{}

func (t proxyTask) EmptyOutput() EmptyStruct {
	return EmptyStruct{}
}

func (t proxyTask) Call(ctx *Context, input EmptyStruct) (EmptyStruct, error) {
	err := t.f(ctx)
	return EmptyStruct{}, err
}

func (t proxyTask) InputType() reflect.Type {
	return reflect.TypeOf(EmptyStruct{})
}
