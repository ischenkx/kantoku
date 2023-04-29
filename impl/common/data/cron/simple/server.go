package simple

import (
	"context"
	"github.com/go-co-op/gocron"
	"kantoku/common/data/pool"
	"log"
)

type Server struct {
	inputs    pool.Reader[Event]
	outputs   pool.Writer[string]
	scheduler *gocron.Scheduler
}

func NewServer(scheduler *gocron.Scheduler, inputs pool.Reader[Event], outputs pool.Writer[string]) *Server {
	return &Server{
		inputs:    inputs,
		outputs:   outputs,
		scheduler: scheduler,
	}
}

func (server *Server) Run(ctx context.Context) error {
	server.scheduler.StartAsync()
	defer server.scheduler.Stop()

	events, err := server.inputs.Read(ctx)
	if err != nil {
		return err
	}

runner:
	for {
		select {
		case <-ctx.Done():
			break runner
		case ev := <-events:
			log.Println("Scheduling:", ev.ID)
			_, err := server.scheduler.
				Every(1).
				Day().
				At(ev.When).
				LimitRunsTo(1).
				Do(func(ctx context.Context, id string) {
					log.Println("Ready:", id)
					if err := server.outputs.Write(ctx, id); err != nil {
						log.Println("failed to write to outputs:", err)
					}
				}, ctx, ev.ID)

			if err != nil {
				log.Println("failed to schedule a job:", err)
			}
		}
	}

	return nil
}
