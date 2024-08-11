package fn

import (
	"errors"
	"reflect"
)

type AbstractFunction[Input, Output any] interface {
	EmptyOutput() Output
	Call(*Context, Input) (Output, error)
	InputType() reflect.Type
	Sched(ctx *Context, input Input) Output
}

type Function[T AbstractFunction[Input, Output], Input, Output any] struct {
	ID string
}

func (f Function[T, Input, Output]) EmptyOutput() Output {
	var zero Output
	return zero
}

func (f Function[T, Input, Output]) InputType() reflect.Type {
	var zero Input
	return reflect.TypeOf(zero)
}

func (f Function[T, Input, Output]) Call(context *Context, input Input) (Output, error) {
	return f.EmptyOutput(), errors.New("not callable")
}

func (f Function[T, Input, Output]) Sched(ctx *Context, input Input) Output {
	out := f.EmptyOutput()
	ctx.Scheduled = append(ctx.Scheduled, ScheduledTask{
		Type:    taskType[T, Input, Output](),
		Inputs:  ctx.addFutureStruct(input, nil),
		Outputs: ctx.addFutureStruct(out, nil),
	})
	return out
}
