package depot

import (
	"context"
	"kantoku"
	"kantoku/common/data/kv"
	"kantoku/common/data/pool"
	"kantoku/common/deps"
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

func (depot *Depot) Deps() deps.Deps {
	return depot.deps
}

func (depot *Depot) Group2Task() kv.Database[string] {
	return depot.group2task
}

func (depot *Depot) Schedule(ctx *kantoku.Context) error {
	rawDependencies, _ := ctx.Data().Get("dependencies")

	dependencies, ok := rawDependencies.([]string)
	if !ok {
		dependencies = nil
	}

	group, err := depot.Deps().MakeGroup(ctx, dependencies...)
	if err != nil {
		return err
	}

	// TODO: possible inconsistency
	if _, err := depot.group2task.Set(ctx, group, ctx.Task.ID()); err != nil {
		return err
	}

	return nil
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
