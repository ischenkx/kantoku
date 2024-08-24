package exe

import (
	"context"
	"errors"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
)

type Router struct {
	typ2exe map[string]executor.Executor
}

func NewRouter() *Router {
	return &Router{
		typ2exe: map[string]executor.Executor{},
	}
}

func (r *Router) AddExecutor(exe executor.Executor, typ string) {
	r.typ2exe[typ] = exe
}

func (r *Router) Execute(ctx context.Context, sys system.AbstractSystem, task task.Task) error {
	typ, ok := task.Info["type"]
	if !ok {
		return errors.New("task without type")
	}
	typStr, ok := typ.(string)
	if !ok {
		return errors.New("task type is not a string")
	}
	exe, ok := r.typ2exe[typStr]
	if !ok {
		return fmt.Errorf("executor for type %s not found", typStr)
	}
	return exe.Execute(ctx, sys, task)
}
