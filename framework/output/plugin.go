package output

import (
	"context"
	"kantoku"
	"kantoku/common/data/kv"
	"kantoku/core/task"
)

type Plugin struct {
	db kv.Database[string, task.Result]
}

func NewPlugin(db kv.Database[string, task.Result]) Plugin {
	return Plugin{db: db}
}

func (p Plugin) Initialize(kantoku *kantoku.Kantoku) {
	kantoku.Props().Set(ResultEvaluator{db: p.db}, "result")
	kantoku.Props().Set(OutputEvaluator{db: p.db}, "output")
}

func (p Plugin) BeforeInitialized(ctx *kantoku.Context) error { return nil }

func (p Plugin) AfterInitialized(ctx *kantoku.Context) {}

func (p Plugin) BeforeScheduled(ctx *kantoku.Context) error {
	return nil
}

func (p Plugin) AfterScheduled(ctx *kantoku.Context) {}

type ResultEvaluator struct {
	db kv.Database[string, task.Result]
}

func (e ResultEvaluator) Evaluate(ctx context.Context, task string) (any, error) {
	return e.db.Get(ctx, task)
}

type OutputEvaluator struct {
	db kv.Database[string, task.Result]
}

func (e OutputEvaluator) Evaluate(ctx context.Context, task string) (any, error) {
	result, err := e.db.Get(ctx, task)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}
