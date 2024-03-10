package resources

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/transport/queue"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"log/slog"
)

type Notifier struct {
	Logger *slog.Logger
	Dst    queue.Publisher[string]
	Topic  string
	DummyObserver
}

func (notifier Notifier) AfterInit(ctx context.Context, resources []resource.Resource) {
	for _, res := range resources {
		fmt.Println("notifier resolution:", res.ID)
		if err := notifier.Dst.Publish(ctx, res.ID); err != nil {
			notifier.Logger.Error("failed to send an initialized resource",
				slog.String("id", res.ID),
				slog.String("error", err.Error()))
		}
	}
}
