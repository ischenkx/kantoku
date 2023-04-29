package chutil

import "context"

func CloseWithContext[Item any](ctx context.Context, ch chan<- Item) {
	go func(ctx context.Context, ch chan<- Item) {
		<-ctx.Done()
		close(ch)
	}(ctx, ch)
}
