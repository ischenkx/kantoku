package status

import (
	"context"
	"kantoku"
	"kantoku/common/data/kv"
)

type Plugin struct {
	db kv.Database[string, Status]
}

func NewPlugin(db kv.Database[string, Status]) *Plugin {
	return &Plugin{db: db}
}

func (p *Plugin) Initialize(kantoku *kantoku.Kantoku) {
	kantoku.Props().Set(Evaluator{db: p.db}, "status")
}

func (p *Plugin) BeforeInitialized(ctx *kantoku.Context) error {
	return nil
}

func (p *Plugin) AfterInitialized(ctx *kantoku.Context) {}

func (p *Plugin) BeforeScheduled(ctx *kantoku.Context) error {
	return nil
}

func (p *Plugin) AfterScheduled(ctx *kantoku.Context) {}

type Evaluator struct {
	db kv.Database[string, Status]
}

func (e Evaluator) Evaluate(ctx context.Context, task string) (any, error) {
	return e.db.Get(ctx, task)
}
