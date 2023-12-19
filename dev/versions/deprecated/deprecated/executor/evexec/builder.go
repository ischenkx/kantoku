package evexec

import (
	"kantoku/backend/executor"
	job2 "kantoku/framework/job"
	"kantoku/job"
)

// Builder is a cosmetic structure to create executor with named args
type Builder[Task job2.Job] struct {
	Runner   executor.Runner[Task, []byte]
	Platform job.Platform[Task]
	Resolver TopicResolver
}

func (builder Builder[Task]) Build() *Executor[Task] {
	return &Executor[Task]{
		runner:        builder.Runner,
		platform:      builder.Platform,
		topicResolver: builder.Resolver,
	}
}
