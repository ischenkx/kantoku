package util

import "context"

func Wait(ctx context.Context) {
	<-ctx.Done()
}
