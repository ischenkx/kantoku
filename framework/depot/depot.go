package depot

import (
	"context"
	"kantoku/common/data/kv"
	"kantoku/common/deps"
	"kantoku/common/pool"
	"log"
)

type Depot struct {
	deps       deps.Deps
	group2task kv.Database[string]
}

func New(deps deps.Deps, group2task kv.Database[string]) *Depot {
	return &Depot{
		deps:       deps,
		group2task: group2task,
	}
}

func (depot *Depot) Schedule(ctx context.Context, id string, dependencies []string) error {
	group, err := depot.Deps().Make(ctx, dependencies...)
	if err != nil {
		return err
	}

	// TODO: possible inconsistency
	if _, err := depot.group2task.Set(ctx, group, id); err != nil {
		return err
	}

	return nil
}

func (depot *Depot) Deps() deps.Deps {
	return depot.deps
}

func (depot *Depot) Process(ctx context.Context, outputs pool.Writer[string]) error {
	ready, err := depot.Deps().Ready(ctx)
	if err != nil {
		return err
	}

loop:
	for {
		log.Println("reading from queue")
		select {
		case <-ctx.Done():
			break loop
		case id := <-ready:
			log.Println("ready:", id)
			taskID, err := depot.group2task.Get(ctx, id)
			if err != nil {
				log.Println("failed to get a task assigned to the group:", err)
				continue
			}
			if err := outputs.Write(ctx, taskID); err != nil {
				log.Println("failed to schedule a task:", err)
				continue
			}
		}
	}

	return nil
}
