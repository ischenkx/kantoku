package arg

import (
	"bytes"
	"context"
	"kantoku/common/codec"
	"kantoku/common/util"
	"kantoku/core/task"
)

type Executor[T task.AbstractTask] struct {
	codec  codec.Codec[[]Arg]
	runner Runner
}

func (executor *Executor[T]) Execute(ctx context.Context, t T) (task.Result, error) {
	raw := t.Argument(ctx)
	args, err := executor.codec.Decode(bytes.NewReader(raw))
	if err != nil {
		return util.Default[task.Result](), err
	}
	ctx = context.WithValue(ctx, "task", t)
	return executor.runner.Run(ctx, t.Type(ctx), args)
}
