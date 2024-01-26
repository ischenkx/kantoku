package functional

import (
	"fmt"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/future"
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

type Task[Input, Output any] interface {
	EmptyOutput() Output
	Call(*Context, Input) (Output, error)
	InputType() reflect.Type
}

type ExampleInput struct {
	a future.Future[int]
	b future.Future[string]
}

type ExampleOutput struct {
	x future.Future[float32]
}

type ExampleTask struct{}

var _ Task[ExampleInput, ExampleOutput] = &ExampleTask{}

func (wt ExampleTask) Call(ctx *Context, input ExampleInput) (ExampleOutput, error) {
	log.Println(input.a)
	log.Println(input.b)
	o := ExampleOutput{x: future.FromValue[float32](4.2)}
	return o, nil
}

func (wt ExampleTask) CallWithSub(ctx *Context, input ExampleInput) (ExampleOutput, error) {
	log.Println(input.a)
	log.Println(input.b)
	promisedOutput := Execute[ExampleTask, ExampleInput, ExampleOutput](ctx, wt, ExampleInput{
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
