package evexec

import (
	"kantoku/backend/executor/common"
	platform2 "kantoku/kernel/platform"
)

// Builder is a cosmetic structure to create executor with named args
type Builder[Task platform2.Task] struct {
	Runner   common.Runner[Task, []byte]
	Platform platform2.Platform[Task]
	Resolver TopicResolver
}

func (builder Builder[Task]) Build() *Executor[Task] {
	return &Executor[Task]{
		runner:        builder.Runner,
		platform:      builder.Platform,
		topicResolver: builder.Resolver,
	}
}
