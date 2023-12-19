package alpha

import "context"

type Plugin = any

type OnNewPlugin interface {
	OnNew(ctx context.Context, alpha Alpha)
}

type OnBeforeRunPlugin interface {
	OnBeforeRun(ctx context.Context, alpha Alpha) error
}

type OnRunPlugin interface {
	OnRun(ctx context.Context, alpha Alpha)
}

type OnReceivePlugin interface {
	OnReceive(ctx context.Context, alpha Alpha) error
}

type OnExecutedPlugin interface {
	OnExecuted(ctx context.Context, alpha Alpha, result Result)
}

type OnResultSaveFailurePlugin interface {
	OnResultSaveFailure(ctx context.Context, alpha Alpha, result Result, err error)
}
