package evexec

import (
	"kantoku/platform"
	"kantoku/unused/backend/executor/common"
)

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
