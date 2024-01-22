package discovery

import (
	"context"
	"errors"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"golang.org/x/sync/errgroup"
)

type Middleware[Service service.Service] struct {
	InfoProvider  func(ctx context.Context, srvc Service) (map[string]any, error)
	Events        *event.Broker
	RequestCodec  codec.Codec[Request, []byte]
	ResponseCodec codec.Codec[Response, []byte]
}

func (m Middleware[Service]) BeforeRun(ctx context.Context, g *errgroup.Group, service service.Service) {
	fmt.Println("1111")

	g.Go(func() error {
		typedService, ok := service.(Service)
		if !ok {
			return errors.New("failed to cast service to the specified type")
		}

		responder := &Responder[Service]{
			Service:       typedService,
			InfoProvider:  m.InfoProvider,
			Events:        m.Events,
			RequestCodec:  m.RequestCodec,
			ResponseCodec: m.ResponseCodec,
		}
		if err := responder.Run(ctx); err != nil {
			return fmt.Errorf("failed to start a service discovery responder: %w", err)
		}

		return nil
	})
}

func WithStaticInfo[Service service.Service](
	info map[string]any,
	events *event.Broker,
	requestCodec codec.Codec[Request, []byte],
	responseCodec codec.Codec[Response, []byte],
) Middleware[Service] {
	return Middleware[Service]{
		InfoProvider: func(ctx context.Context, srvc Service) (map[string]any, error) {
			return info, nil
		},
		Events:        events,
		RequestCodec:  requestCodec,
		ResponseCodec: responseCodec,
	}
}
