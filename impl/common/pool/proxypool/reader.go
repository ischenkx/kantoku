package proxypool

import (
	"context"
	"kantoku/common/chutil"
	"kantoku/common/data/pool"
	transformator2 "kantoku/common/transformer"
)

type Reader[In, Out any] struct {
	inputs        pool.Reader[In]
	transformator transformator2.Transformer[In, Out]
}

func NewReader[In, Out any](inputs pool.Reader[In], transformator transformator2.Transformer[In, Out]) *Reader[In, Out] {
	return &Reader[In, Out]{
		inputs:        inputs,
		transformator: transformator,
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
			outputs <- r.transformator(item)
		}
	}(ctx, inputs, newInputs)

	return newInputs, nil
}
