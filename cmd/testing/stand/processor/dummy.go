package main

import (
	"github.com/ischenkx/kantoku/pkg/lib/tasks/functional"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/future"
)

type DummyTask struct{}

func (t DummyTask) Call(context *functional.Context, input MathInput) (MathOutput, error) {
	return functional.Execute[AddTask, MathInput, MathOutput](context, AddTask{}, MathInput{
		Left:  future.FromValue[int](123),
		Right: future.FromValue[int](123),
	}), nil
}
