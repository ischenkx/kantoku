package kantoku

import (
	"context"
	"fmt"
	"kantoku/framework/argument"
)

type Option func(ctx *Context)

type Spec struct {
	Type      string
	Arguments []any
	Options   []Option
}

func Task(typ string, args ...any) Spec {
	return Spec{
		Type:      typ,
		Arguments: args,
	}
}

func (spec Spec) With(options ...Option) Spec {
	spec.Options = append(spec.Options, options...)
	return spec
}

type TaskInstance struct {
	id        string
	typ       string
	arguments []any
}

func (instance *TaskInstance) ID() string {
	return instance.id
}

func (instance *TaskInstance) Type() string {
	return instance.typ
}

func (instance *TaskInstance) Arg(index int) (any, bool) {
	if index < 0 || len(instance.arguments) <= index {
		return nil, false
	}

	return instance.arguments[index], true
}

func (instance *TaskInstance) CountArgs() int {
	return len(instance.arguments)
}

type StoredTask struct {
	ID        string
	Type      string
	Arguments []argument.Argument
}

type ScheduledTask struct {
	id string
}

func (s ScheduledTask) ID() string {
	return s.id
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

func (view *TaskView) Arguments(ctx context.Context) ([]argument.Argument, error) {
	if err := view.loadStored(ctx); err != nil {
		return nil, err
	}
	return view.stored.Arguments, nil
}

func (view *TaskView) Prop(ctx context.Context, path ...string) (any, error) {
	evaluator, ok := view.Kantoku().Props().Get(path...)
	if !ok {
		return nil, fmt.Errorf("no evaluator provided for key: %s", path)
	}

	return evaluator.Evaluate(ctx, view.id)
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
