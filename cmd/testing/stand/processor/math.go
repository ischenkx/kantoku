package main

import (
	"github.com/ischenkx/kantoku/pkg/lib/tasks/functional"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/future"
	"reflect"
)

type MathInput struct {
	Left  future.Future[int]
	Right future.Future[int]
}

type MathOutput struct {
	Result future.Future[int]
}

type AddTask struct{ MathTask }

func (t AddTask) Call(context *functional.Context, input MathInput) (MathOutput, error) {
	return MathOutput{
		Result: future.FromValue(input.Left.Value() + input.Right.Value()),
	}, nil
}

type MulTask struct{ MathTask }

func (t MulTask) Call(context functional.Context, input MathInput) (MathOutput, error) {
	return MathOutput{
		Result: future.FromValue(input.Left.Value() * input.Right.Value()),
	}, nil
}

type DivTask struct{ MathTask }

func (t DivTask) Call(context functional.Context, input MathInput) (MathOutput, error) {
	return MathOutput{
		Result: future.FromValue(input.Left.Value() / input.Right.Value()),
	}, nil
}

type MathTask struct{}

func (t MathTask) EmptyOutput() MathOutput {
	return MathOutput{Result: future.Empty[int]()}
}

func (t MathTask) InputType() reflect.Type {
	return reflect.TypeOf(MathInput{})
}

type SumInput struct {
	Args future.Future[[]int]
}

type SumTask struct{ MathTask }

func (t SumTask) Call(ctx *functional.Context, input SumInput) (MathOutput, error) {
	if len(input.Args.Value()) == 0 {
		return MathOutput{Result: future.FromValue(0)}, nil
	}
	if len(input.Args.Value()) == 1 {
		return MathOutput{Result: future.FromValue(input.Args.Value()[0])}, nil
	}
	sumSuffix := functional.Execute[SumTask, SumInput, MathOutput](ctx, SumTask{}, SumInput{
		Args: future.FromValue(input.Args.Value()[1:]),
	})
	total := functional.Execute[AddTask, MathInput, MathOutput](ctx, AddTask{}, MathInput{
		Left:  future.FromValue(input.Args.Value()[0]),
		Right: sumSuffix.Result,
	})
	return MathOutput{Result: total.Result}, nil
}

func (t SumTask) InputType() reflect.Type {
	return reflect.TypeOf(SumInput{})
}
