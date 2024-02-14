package resources

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"log/slog"
)

type Notifier struct {
	Logger *slog.Logger
	Broker *event.Broker
	Topic  string
	DummyObserver
}

func (notifier Notifier) AfterInit(ctx context.Context, resources []resource.Resource) {
	for _, res := range resources {
		ev := event.New(notifier.Topic, []byte(res.ID))
		if err := notifier.Broker.Send(ctx, ev); err != nil {
			notifier.Logger.Error("failed to send an initialized resource",
				slog.String("id", res.ID),
				slog.String("error", err.Error()))
		}
	}
}
