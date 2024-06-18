package discovery

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"log/slog"
)

const (
	RequestsTopic  = "discovery:request"
	ResponsesTopic = "discovery:response"
)

type Request struct {
	ID string
}

type Response struct {
	RequestID   string
	ServiceInfo ServiceInfo
}

type Responder[Service service.Service] struct {
	Service       Service
	InfoProvider  func(ctx context.Context, srvc Service) (map[string]any, error)
	Events        *event.Broker
	RequestCodec  codec.Codec[Request, []byte]
	ResponseCodec codec.Codec[Response, []byte]
}

func (responder *Responder[Service]) Run(ctx context.Context) error {
	channel, err := responder.Events.Consume(ctx, broker.TopicsInfo{
		Group: responder.Service.ID(),
		Topics: []string{
			RequestsTopic,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	//responder.Service.Logger().Info("starting a responder")

	broker.Processor[event.Event]{
		Handler: func(ctx context.Context, ev event.Event) error {
			request, err := responder.RequestCodec.Decode(ev.Data)
			if err != nil {
				responder.Service.Logger().Error("failed to decode a discovery request",
					slog.String("error", err.Error()),
					slog.String("event_id", ev.ID))
				return nil
			}

			//responder.Service.Logger().Info("received a discovery request",
			//	slog.String("id", request.ID))

			var info map[string]any
			if responder.InfoProvider != nil {
				info, err = responder.InfoProvider(ctx, responder.Service)
				if err != nil {
					responder.Service.Logger().Error("failed to get service info",
						slog.String("error", err.Error()),
						slog.String("event_id", ev.ID))
					return nil
				}
			}

			response := Response{
				RequestID: request.ID,
				ServiceInfo: ServiceInfo{
					ID:   responder.Service.ID(),
					Name: responder.Service.Name(),
					Info: info,
				},
			}

			encodedResponse, err := responder.ResponseCodec.Encode(response)
			if err != nil {
				responder.Service.Logger().Error("failed to encode a discovery response",
					slog.String("error", err.Error()),
					slog.String("event_id", ev.ID))
				return nil
			}

			if err := responder.Events.Send(ctx, event.New(ResponsesTopic, encodedResponse)); err != nil {
				responder.Service.Logger().Error("failed to send a discovery response",
					slog.String("error", err.Error()),
					slog.String("event_id", ev.ID))
				return nil
			}

			return nil
		},
	}.Process(ctx, channel)

	return nil
}
