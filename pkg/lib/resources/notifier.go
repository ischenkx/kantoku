package resources

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/core"
	"log/slog"
)

type Notifier struct {
	Broker core.Broker
	Topic  string

	Logger *slog.Logger

	DummyObserver
}

func (notifier Notifier) AfterInit(ctx context.Context, resources []core.Resource) {
	for _, res := range resources {
		ev := core.NewEvent(notifier.Topic, []byte(res.ID))
		if err := notifier.Broker.Send(ctx, ev); err != nil {
			notifier.Logger.Error("failed to send an initialized resource_db",
				slog.String("id", res.ID),
				slog.String("error", err.Error()))
		}
	}
}
