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
	Call(Context, Input) (Output, error)
	InputType() reflect.Type
}

type MyInput struct {
	a future.Future[int]
	b future.Future[string]
}

type MyOutput struct {
	x future.Future[float32]
}

type WorkTask struct{}

var _ Task[MyInput, MyOutput] = &WorkTask{}

func (wt WorkTask) Call(ctx Context, input MyInput) (MyOutput, error) {
	log.Println(input.a)
	log.Println(input.b)
	o := MyOutput{x: future.FromValue[float32](4.2)}
	return o, nil
}

func (wt WorkTask) CallWithSub(ctx Context, input MyInput) (MyOutput, error) {
	log.Println(input.a)
	log.Println(input.b)
	promisedOutput := Execute[WorkTask, MyInput, MyOutput](ctx, wt, MyInput{
		a: future.FromValue[int](len(input.b.Value())),
		b: future.FromValue[string](fmt.Sprint(input.a.Value())),
	})
	return promisedOutput, nil
}

func (wt WorkTask) EmptyFutureInput() MyInput {
	return MyInput{
		a: future.Future[int]{},
		b: future.Future[string]{},
	}
}

func (wt WorkTask) EmptyOutput() MyOutput {
	return MyOutput{
		x: future.Future[float32]{},
	}
}

func (wt WorkTask) InputType() reflect.Type {
	return reflect.TypeOf(MyInput{})
}
