package evexec

import (
	"kantoku/backend/executor/common"
	"kantoku/kernel/platform"
)

// Builder is a cosmetic structure to create executor with named args
type Builder[Task platform.Task] struct {
	Runner   common.Runner[Task, []byte]
	Platform platform.Platform[Task]
	Resolver TopicResolver
}

func (builder Builder[Task]) Build() *Executor[Task] {
	return &Executor[Task]{
		runner:        builder.Runner,
		platform:      builder.Platform,
		topicResolver: builder.Resolver,
	}
}
