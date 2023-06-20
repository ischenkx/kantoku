package rider

import "context"

type Orchestrator interface {
	Orchestrate(ctx context.Context, jobs ...Job) error
}
