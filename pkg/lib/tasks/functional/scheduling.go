package functional

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"reflect"
)

func SchedulingContext(ctx context.Context, sys system.AbstractSystem, f func(*Context) error) error {
	proxy := proxyTask{f: f}
	exe := NewExecutor[proxyTask, EmptyStruct, EmptyStruct](proxy)

	sysTask := task.Task{
		Inputs:  []resource.ID{},
		Outputs: []resource.ID{},
		ID:      "proxy---should-not-see-this",
		Info:    record.R{},
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
