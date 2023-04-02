package proxypool

import (
	"context"
	"kantoku/common/chutil"
	"kantoku/common/data/pool"
)

type Reader[In, Out any] struct {
	inputs      pool.Reader[In]
	transformer func(ctx context.Context, input In) (Out, bool)
}

func NewReader[In, Out any](inputs pool.Reader[In], transformer func(ctx context.Context, input In) (Out, bool)) *Reader[In, Out] {
	return &Reader[In, Out]{
		inputs:      inputs,
		transformer: transformer,
	}
}

func (r *Reader[In, Out]) Read(ctx context.Context) (<-chan Out, error) {
	inputs, err := r.inputs.Read(ctx)
	if err != nil {
		return nil, err
	}

	newInputs := make(chan Out, 1024)

	chutil.SyncWithContext(ctx, newInputs)
	go func(ctx context.Context, inputs <-chan In, outputs chan<- Out) {
		for item := range inputs {
			output, ok := r.transformer(ctx, item)
			if ok {
				outputs <- output
			}
		}
	}(ctx, inputs, newInputs)

	return newInputs, nil
}
