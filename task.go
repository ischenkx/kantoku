package kantoku

import (
	"context"
	"fmt"
	"kantoku/platform"
)

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

type TaskInstance struct {
	id   string
	typ  string
	data []byte
}

func (instance *TaskInstance) ID() string {
	return instance.id
}

func (instance *TaskInstance) Type() string {
	return instance.typ
}

func (instance *TaskInstance) Data() []byte {
	return instance.data
}

type StoredTask struct {
	Id   string
	Type string
	Data []byte
}

func (task StoredTask) ID() string {
	return task.Id
}

type View struct {
	kantoku *Kantoku
	id      string
	stored  *StoredTask
}

func (view *View) Kantoku() *Kantoku {
	return view.kantoku
}

func (view *View) ID() string {
	return view.id
}

func (view *View) Type(ctx context.Context) (string, error) {
	stored, err := view.Stored(ctx)
	return stored.Type, err
}

func (view *View) Data(ctx context.Context) ([]byte, error) {
	stored, err := view.Stored(ctx)
	return stored.Data, err
}

func (view *View) Prop(ctx context.Context, path ...string) (any, error) {
	evaluator, ok := view.Kantoku().Props().Get(path...)
	if !ok {
		return nil, fmt.Errorf("no evaluator provided for key: %s", path)
	}

	return evaluator.Evaluate(ctx, view.id)
}

func (view *View) Stored(ctx context.Context) (TaskInstance, error) {
	if view.stored != nil {
		return *view.stored, nil
	}

	stored, err := view.kantoku.tasks.Get(ctx, view.id)
	if err != nil {
		return TaskInstance{}, err
	}
	view.stored = &stored

	return stored, nil
}

func (view *View) Result(ctx context.Context) (platform.Result, error) {
	return view.Kantoku().Outputs().Get(ctx, view.ID())
}
