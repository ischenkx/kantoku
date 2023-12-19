package job

import (
	"context"
)

type Option func(ctx *Context) error

// Spec is an abstract representation of a task
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

// Job is a compiled spec.
type Job struct {
	Id   string
	Type string
	Data []byte
}

func (job Job) ID() string {
	return job.Id
}

// View is a helper structure that provides convenient methods to work
// with a task
type View struct {
	kernel   *Manager
	instance *Job
	id       string
}

func (view *View) Kernel() *Manager {
	return view.kernel
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

func (view *View) Instance(ctx context.Context) (Job, error) {
	if view.instance != nil {
		return *view.instance, nil
	}

	instance, err := view.kernel.Tasks().Get(ctx, view.id)
	if err != nil {
		return Job{}, err
	}
	view.instance = &instance

	return instance, nil
}

func (view *View) Result(ctx context.Context) (Result, error) {
	return view.Kernel().Outputs().Get(ctx, view.ID())
}
