package discovery

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/data/uid"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"github.com/ischenkx/kantoku/pkg/core"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"time"
)

type Poller struct {
	Hub           Hub
	Events        core.Broker
	RequestCodec  codec.Codec[Request, []byte]
	ResponseCodec codec.Codec[Response, []byte]
	Interval      time.Duration

	service.Core
}

func (poller *Poller) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		poller.poll(ctx)
		return nil
	})

	g.Go(func() error {
		return poller.collectResponses(ctx)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (poller *Poller) poll(ctx context.Context) {
	ticker := time.NewTicker(poller.Interval)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			request := Request{ID: uid.Generate()}
			encodedRequest, err := poller.RequestCodec.Encode(request)
			if err != nil {
				poller.Logger().Error("failed to encode a discovery request",
					slog.String("error", err.Error()))
				continue
			}

			if err := poller.Events.Send(ctx, core.NewEvent(RequestsTopic, encodedRequest)); err != nil {
				poller.Logger().Error("failed to send a discovery request",
					slog.String("error", err.Error()))
				continue
			}

			poller.Logger().Info("sent a poll request")
		}
	}
}

func (poller *Poller) collectResponses(ctx context.Context) error {
	channel, err := poller.Events.Consume(ctx,
		[]string{ResponsesTopic},
		broker.ConsumerSettings{
			Group:                poller.Core.Name(),
			InitializationPolicy: broker.NewestOffset,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to consume responses: %w", err)
	}

	broker.Processor[core.Event]{
		Handler: func(ctx context.Context, ev core.Event) error {
			response, err := poller.ResponseCodec.Decode(ev.Data)
			if err != nil {
				poller.Logger().Error("failed to decode a discovery response",
					slog.String("error", err.Error()),
					slog.String("event_id", ev.ID))
				return nil
			}

			poller.Logger().Info("received a discovery response",
				slog.String("event_id", ev.ID),
				slog.String("request_id", response.RequestID),
				slog.String("service.name", response.ServiceInfo.Name),
				slog.String("service.id", response.ServiceInfo.ID))

			if err := poller.Hub.Register(ctx, response.ServiceInfo); err != nil {
				poller.Logger().Error("failed to register a service",
					slog.String("error", err.Error()))
				return nil
			}

			return nil
		},
	}.Process(ctx, channel)

	return nil
}
