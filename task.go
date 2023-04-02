package kantoku

import (
	"context"
	"fmt"
)

type Option func(ctx *Context)

type Spec struct {
	Type    string
	Data    any
	Options []Option
}

func Task(typ string, data any) Spec {
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
	data any
}

func (instance *TaskInstance) ID() string {
	return instance.id
}

func (instance *TaskInstance) Type() string {
	return instance.typ
}

func (instance *TaskInstance) Data() any {
	return instance.data
}

type StoredTask struct {
	Id   string
	Type string
	Data any
}

func (task StoredTask) ID() string {
	return task.Id
}

type TaskView struct {
	kantoku *Kantoku
	id      string
	stored  *StoredTask
}

func (view *TaskView) Kantoku() *Kantoku {
	return view.kantoku
}

func (view *TaskView) ID() string {
	return view.id
}

func (view *TaskView) Type(ctx context.Context) (string, error) {
	if err := view.loadStored(ctx); err != nil {
		return "", err
	}
	return view.stored.Type, nil
}

func (view *TaskView) Data(ctx context.Context) (any, error) {
	if err := view.loadStored(ctx); err != nil {
		return nil, err
	}
	return view.stored.Data, nil
}

func (view *TaskView) Prop(ctx context.Context, path ...string) (any, error) {
	evaluator, ok := view.Kantoku().Props().Get(path...)
	if !ok {
		return nil, fmt.Errorf("no evaluator provided for key: %s", path)
	}

	return evaluator.Evaluate(ctx, view.id)
}

func (view *TaskView) AsStored(ctx context.Context) (StoredTask, error) {
	err := view.loadStored(ctx)
	if err != nil {
		return StoredTask{}, err
	}
	return *view.stored, nil
}

func (view *TaskView) loadStored(ctx context.Context) error {
	if view.stored != nil {
		return nil
	}

	task, err := view.Kantoku().tasks.Get(ctx, view.id)
	if err != nil {
		return err
	}
	view.stored = &task

	return nil
}
