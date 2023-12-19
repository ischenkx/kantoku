package roller

import "context"

type Action struct {
	Func     func(ctx context.Context)
	Rollback func(ctx context.Context)
	Label    string
}
