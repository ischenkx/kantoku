package depot

import (
	"context"
	"kantoku"
	"kantoku/common/data/bimap"
	"kantoku/common/data/pool"
	"kantoku/common/deps"
	"log"
)

type Depot struct {
	deps           deps.Deps
	groupTaskBimap bimap.Bimap[string, string]
}

func New(deps deps.Deps, groupTaskBimap bimap.Bimap[string, string]) *Depot {
	return &Depot{
		deps:           deps,
		groupTaskBimap: groupTaskBimap,
	}
}

func (depot *Depot) Deps() deps.Deps {
	return depot.deps
}

func (depot *Depot) GroupTaskBimap() bimap.Bimap[string, string] {
	return depot.groupTaskBimap
}

func (depot *Depot) Write(ctx context.Context, task *kantoku.TaskInstance) error {
	data := kantoku.GetPluginData(ctx).GetWithDefault("dependencies", &PluginData{}).(*PluginData)

	group, err := depot.Deps().MakeGroup(ctx, data.Dependencies...)
	if err != nil {
		return err
	}

	// TODO: possible inconsistency
	if err := depot.groupTaskBimap.Save(ctx, group, task.ID); err != nil {
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
			taskID, err := depot.groupTaskBimap.ByKey(ctx, id)
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
