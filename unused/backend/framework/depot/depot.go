package depot

import (
	"context"
	"kantoku"
	"kantoku/common/data/bimap"
	"kantoku/platform"
	"kantoku/unused/backend/framework/depot/deps"
	"log"
)

type Depot struct {
	deps           deps.Deps
	groupTaskBimap bimap.Bimap[string, string]
	inputs         platform.Inputs
}

func New(deps deps.Deps, groupTaskBimap bimap.Bimap[string, string], inputs platform.Inputs) *Depot {
	return &Depot{
		deps:           deps,
		groupTaskBimap: groupTaskBimap,
		inputs:         inputs,
	}
}

func (depot *Depot) Deps() deps.Deps {
	return depot.deps
}

func (depot *Depot) GroupTaskBimap() bimap.Bimap[string, string] {
	return depot.groupTaskBimap
}

func (depot *Depot) Write(ctx context.Context, id string) error {
	data := kantoku.GetPluginData(ctx).GetWithDefault("dependencies", &PluginData{}).(*PluginData)

	group, err := depot.Deps().MakeGroup(ctx, data.Dependencies...)
	if err != nil {
		return err
	}

	// TODO: possible inconsistency
	if err := depot.groupTaskBimap.Save(ctx, group, id); err != nil {
		return err
	}

	return nil
}

func (depot *Depot) Read(ctx context.Context) (<-chan string, error) {
	return depot.inputs.Read(ctx)
}

func (depot *Depot) Process(ctx context.Context) error {
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
			if err := depot.inputs.Write(ctx, taskID); err != nil {
				log.Println("failed to schedule a task:", err)
				continue
			}
		}
	}

	return nil
}
