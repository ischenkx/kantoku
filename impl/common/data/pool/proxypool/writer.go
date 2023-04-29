package proxypool

import (
	"context"
	"kantoku/common/data/pool"
	transformer "kantoku/common/transformer"
)

type Writer[In, Out any] struct {
	outputs       pool.Writer[Out]
	transformator transformer.Transformer[In, Out]
}

func NewWriter[In, Out any](outputs pool.Writer[Out], transformator transformer.Transformer[In, Out]) *Writer[In, Out] {
	return &Writer[In, Out]{
		outputs:       outputs,
		transformator: transformator,
	}
}

func (w *Writer[In, Out]) Write(ctx context.Context, item In) error {
	return w.outputs.Write(ctx, w.transformator(item))
}
