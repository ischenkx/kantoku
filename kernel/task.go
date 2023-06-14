package kernel

import (
	"context"
	"kantoku/kernel/platform"
)

// Spec is an abstract representation of a task

type Option func(ctx *Context) error

type Spec struct {
	Type    string
	Data    []byte
	Options []Option
}

func Describe(typ string, data []byte) Spec {
	return Spec{
		Type: typ,
		Data: data,
	}
}

func (spec Spec) With(options ...Option) Spec {
	spec.Options = append(spec.Options, options...)
	return spec
}

// Task is a compiled spec.
type Task struct {
	Id   string
	Type string
	Data []byte
}

func (task Task) ID() string {
	return task.Id
}

// View is a helper structure that provides convenient methods to work
// with a task

type View struct {
	kantoku  *Kernel
	id       string
	instance *Task
}

func (view *View) Kantoku() *Kernel {
	return view.kantoku
}

func (view *View) ID() string {
	return view.id
}

func (view *View) Type(ctx context.Context) (string, error) {
	stored, err := view.Instance(ctx)
	return stored.Type, err
}

func (view *View) Data(ctx context.Context) ([]byte, error) {
	stored, err := view.Instance(ctx)
	return stored.Data, err
}

func (view *View) Instance(ctx context.Context) (Task, error) {
	if view.instance != nil {
		return *view.instance, nil
	}

	instance, err := view.kantoku.platform.DB().Get(ctx, view.id)
	if err != nil {
		return Task{}, err
	}
	view.instance = &instance

	return instance, nil
}

func (view *View) Result(ctx context.Context) (platform.Result, error) {
	return view.Kantoku().Outputs().Get(ctx, view.ID())
}
