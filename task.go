package kantoku

import (
	"context"
	"fmt"
	"kantoku/platform"
)

// Spec is an abstract representation of a task

type Option func(ctx *Context) error

type Spec struct {
	Type    string
	Data    []byte
	Options []Option
}

func Task(typ string, data []byte) Spec {
	return Spec{
		Type: typ,
		Data: data,
	}
}

func (spec Spec) With(options ...Option) Spec {
	spec.Options = append(spec.Options, options...)
	return spec
}

// TaskInstance is a compiled spec.
type TaskInstance struct {
	Id   string
	Type string
	Data []byte
}

func (instance TaskInstance) ID() string {
	return instance.Id
}

// View is a helper structure that provides convenient methods to work
// with a task

type View struct {
	kantoku  *Kantoku
	id       string
	instance *TaskInstance
}

func (view *View) Kantoku() *Kantoku {
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

func (view *View) Prop(ctx context.Context, path ...string) (any, error) {
	evaluator, ok := view.Kantoku().Props().Get(path...)
	if !ok {
		return nil, fmt.Errorf("no evaluator provided for key: %s", path)
	}

	return evaluator.Evaluate(ctx, view.id)
}

func (view *View) Instance(ctx context.Context) (TaskInstance, error) {
	if view.instance != nil {
		return *view.instance, nil
	}

	instance, err := view.kantoku.platform.DB().Get(ctx, view.id)
	if err != nil {
		return TaskInstance{}, err
	}
	view.instance = &instance

	return instance, nil
}

func (view *View) Result(ctx context.Context) (platform.Result, error) {
	return view.Kantoku().Outputs().Get(ctx, view.ID())
}
