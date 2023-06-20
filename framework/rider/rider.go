package rider

import "context"

type Rider struct {
	jobs         []Job
	orchestrator Orchestrator
}

func (rider *Rider) Register(job Job) {
	rider.jobs = append(rider.jobs, job)
}

func (rider *Rider) List() []Job {
	return rider.jobs
}

func (rider *Rider) Run(ctx context.Context) error {
	return rider.orchestrator.Orchestrate(ctx, rider.List()...)
}
