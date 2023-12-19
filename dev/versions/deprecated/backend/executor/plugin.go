package executor

import (
	"context"
	"kantoku/framework/job"
)

type Plugin any

type ReceivedTaskPlugin interface {
	ReceivedTask(ctx context.Context, id string)
}

type ExecutedTaskPlugin interface {
	ExecutedTask(ctx context.Context, id string, result job.Result)
}

type SavedTaskResultPlugin interface {
	SavedTaskResult(ctx context.Context, id string, result job.Result)
}

type FailedToSaveTaskResultPlugin interface {
	FailedToSaveTaskResult(ctx context.Context, id string, result job.Result)
}
