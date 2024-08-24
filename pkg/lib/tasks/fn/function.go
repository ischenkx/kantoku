package fn

import (
	"errors"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn/future"
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
	EmptyOutput() (Output, error)
	Call(*Context, Input) (Output, error)
	InputType() reflect.Type
	Sched(ctx *Context, input Input) (Output, error)
}

type Function[T AbstractFunction[Input, Output], Input, Output any] struct {
	ID string
}

func (task Function[T, Input, Output]) EmptyOutput() (Output, error) {
	var output Output

	val := reflect.ValueOf(&output)

	numField := val.Elem().NumField()
	for fieldIdx := 0; fieldIdx < numField; fieldIdx++ {
		field := val.Elem().Field(fieldIdx)
		ptr := field.Addr().Interface()
		fut, ok := ptr.(future.InitializeableFuture)
		if !ok {
			return output, fmt.Errorf("failed to convert %v to future.AbstractFuture", ptr)
		}

		fut.Initialize()
	}

	return output, nil
}

func (task Function[T, Input, Output]) InputType() reflect.Type {
	var zero Input
	return reflect.TypeOf(zero)
}

func (task Function[T, Input, Output]) Call(context *Context, input Input) (Output, error) {
	var output Output
	return output, errors.New("not callable")
}

func (task Function[T, Input, Output]) Sched(ctx *Context, input Input) (Output, error) {
	output, err := task.EmptyOutput()
	if err != nil {
		return output, fmt.Errorf("failed to get an empty output: %w", err)
	}

	boundInput, err := ctx.bindObjectToResources(input, nil)
	if err != nil {
		return output, fmt.Errorf("failed to bind inputs: %w", err)
	}

	boundOutput, err := ctx.bindObjectToResources(output, nil)
	if err != nil {
		return output, fmt.Errorf("failed to bind outputs: %w", err)
	}

	ctx.Scheduled = append(ctx.Scheduled, ScheduledTask{
		Type:    taskType[T, Input, Output](),
		Inputs:  boundInput,
		Outputs: boundOutput,
	})

	return output, nil
}

func Func[F AbstractFunction[I, O], I, O any]() F {
	var f F
	return f
}

func Sched[F AbstractFunction[I, O], I, O any](c *Context, input I) (O, error) {
	var f F
	return f.Sched(c, input)
}
