package fn_d

import (
	"errors"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn_d/future"
	"log"
	"reflect"
)

/*
 main ideas:
 task is a function
 it receives input struct and takes data from it for work
 it produces output struct where each field is a Future and may or may not be filled
 (otherwise you can't delegate calculating something to other task)
 for now both input and output ignore everything not serializable and care only about top level fields
*/

/*
ok, getting generic types is kinda hard, so i'll add methods to return type
https://github.com/golang/go/issues/54393
*/

type AbstractFunction[Input, Output any] interface {
	EmptyOutput() Output
	Call(*Context, Input) (Output, error)
	InputType() reflect.Type
	Sched(ctx *Context, input Input) Output
}

type Function[T AbstractFunction[Input, Output], Input, Output any] struct {
	ID string
}

func (task Function[T, Input, Output]) EmptyOutput() Output {
	var zero Output
	return zero
}

func (task Function[T, Input, Output]) InputType() reflect.Type {
	var zero Input
	return reflect.TypeOf(zero)
}

func (task Function[T, Input, Output]) Call(context *Context, input Input) (Output, error) {
	return task.EmptyOutput(), errors.New("not callable")
}

func (task Function[T, Input, Output]) Sched(ctx *Context, input Input) Output {
	out := task.EmptyOutput()
	ctx.Scheduled = append(ctx.Scheduled, ScheduledTask{
		Type:    taskType[T, Input, Output](),
		Inputs:  ctx.addFutureStruct(input, nil),
		Outputs: ctx.addFutureStruct(out, nil),
	})
	return out
}

func Func[F AbstractFunction[I, O], I, O any]() F {
	var f F
	return f
}

type ExampleInput struct {
	a future.Future[int]
	b future.Future[string]
}

type ExampleOutput struct {
	x future.Future[float32]
}

type ExampleTask struct {
	Function[ExampleTask, ExampleInput, ExampleOutput]
}

var _ AbstractFunction[ExampleInput, ExampleOutput] = ExampleTask{}

func (wt ExampleTask) Call(ctx *Context, input ExampleInput) (ExampleOutput, error) {
	log.Println(input.a)
	log.Println(input.b)
	o := ExampleOutput{x: future.FromValue[float32](4.2)}
	return o, nil
}

func (wt ExampleTask) CallWithSub(ctx *Context, input ExampleInput) (ExampleOutput, error) {
	log.Println(input.a)
	log.Println(input.b)
	promisedOutput := wt.Sched(ctx, ExampleInput{
		a: future.FromValue[int](len(input.b.Value())),
		b: future.FromValue[string](fmt.Sprint(input.a.Value())),
	})

	return promisedOutput, nil
}

func (wt ExampleTask) EmptyFutureInput() ExampleInput {
	return ExampleInput{
		a: future.Future[int]{},
		b: future.Future[string]{},
	}
}

func (wt ExampleTask) EmptyOutput() ExampleOutput {
	return ExampleOutput{
		x: future.Future[float32]{},
	}
}

func (wt ExampleTask) InputType() reflect.Type {
	return reflect.TypeOf(ExampleInput{})
}
