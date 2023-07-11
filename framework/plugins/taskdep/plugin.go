package taskdep

import (
	"context"
	"kantoku/framework/infra/demon"
	"kantoku/framework/plugins/depot"
	"kantoku/framework/utils/demons"
	"kantoku/kernel"
	"kantoku/kernel/platform"
	"log"
)

type Plugin struct {
	manager *Manager
	events  platform.Broker
	topic   string
}

func NewPlugin(manager *Manager, events platform.Broker, topic string) *Plugin {
	return &Plugin{
		manager: manager,
		events:  events,
		topic:   topic,
	}
}

type PluginData struct {
	Subtasks []string
}

func (p *Plugin) BeforeScheduled(ctx *kernel.Context) error {
	data := ctx.Data().GetOrSet("taskdep", func() any { return &PluginData{} }).(*PluginData)

	for _, subtask := range data.Subtasks {
		dependency, err := p.manager.SubtaskDependency(ctx, subtask)
		if err != nil {
			return err
		}

		if err := depot.Dependency(dependency)(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (p *Plugin) Demons(ctx context.Context) []demon.Demon {
	return demons.Multi{
		demons.TryProvider(p.manager.deps),
		demons.TryProvider(p.manager.task2dep),
		demons.Functional("TASKDEP_PROCESSOR", p.process),
	}.Demons(ctx)
}

func (p *Plugin) process(ctx context.Context) error {
	listener := p.events.Listen()
	defer listener.Close(ctx)

	err := listener.Subscribe(ctx, p.topic)
	if err != nil {
		return err
	}

	channel, err := listener.Incoming(ctx)
	if err != nil {
		return err
	}
updater:
	for {
		select {
		case <-ctx.Done():
			break updater
		case ev := <-channel:
			id := string(ev.Data)
			if err := p.manager.ResolveTask(ctx, id); err != nil {
				log.Println("failed to resolve a task:", err)
			}
		}
	}

	return nil
}
