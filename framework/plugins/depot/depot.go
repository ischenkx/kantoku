package depot

import (
	"context"
	"fmt"
	"kantoku/common/data/bimap"
	"kantoku/framework/plugins/depot/deps"
	"kantoku/kernel"
	"kantoku/kernel/platform"
	"log"
)

type PluginData struct {
	Dependencies []string
}

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
	data := kernel.GetPluginData(ctx).GetWithDefault("dependencies", &PluginData{}).(*PluginData)

	group, err := depot.Deps().MakeGroup(ctx, data.Dependencies...)
	if err != nil {
		return fmt.Errorf("failed to make a dependency group: %s", err)
	}

	// TODO: possible inconsistency
	// (if the task dependency group had been resolved and processed before the following line was executed)
	if err := depot.groupTaskBimap.Save(ctx, group, id); err != nil {
		return fmt.Errorf("failed to save the (group, task) pair in the bimap: %s", err)
	}

	return nil
}

func (depot *Depot) Read(ctx context.Context) (<-chan string, error) {
	return depot.inputs.Read(ctx)
}

func (depot *Depot) Process(ctx context.Context) error {
	ready, err := depot.Deps().Ready(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize the 'ready' channel: %s", err)
	}

loop:
	for {
		log.Println("reading from queue")
		select {
		case <-ctx.Done():
			break loop
		case id := <-ready:
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
