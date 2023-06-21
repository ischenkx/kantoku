package future

import "context"

type Runner interface {
	Run(ctx context.Context, resolution Resolution)
}

type SequentialRunner []Runner

func (runner SequentialRunner) Run(ctx context.Context, res Resolution) {
	for _, subRunner := range runner {
		subRunner.Run(ctx, res)
	}
}
