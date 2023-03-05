package plugins

import (
	"context"
	"kantoku"
)

// ArgumentCallback is created in plugin and called by kantoku app to create new argument from parameters
// given to the plugin or somehow process the task inside the plugin
type ArgumentCallback interface {
	// Resolve receives the task at the begging of scheduling (and before WithCallback's were resolved)
	Resolve(ctx context.Context, task kantoku.Task, kantoku kantoku.Kantoku) (*kantoku.Argument, error)
	// Register receives the task after it was scheduled (but before WithCallback register calls)
	Register(ctx context.Context, task kantoku.Task, kantoku2 kantoku.Kantoku) error
}

// WithCallback is created in plugin and called by kantoku app after all Arguments have been resolved to
// modify task (add dependencies) or somehow register process the task inside the plugin
type WithCallback interface {
	// Resolve receives the task at the begging of scheduling (but after ArgumentCallback's were resolved)
	Resolve(ctx context.Context, task kantoku.Task, kantoku kantoku.Kantoku) error
	// Register receives the task after it was scheduled (and after ArgumentCallback register calls)
	Register(ctx context.Context, task kantoku.Task, kantoku kantoku.Kantoku) error
}
