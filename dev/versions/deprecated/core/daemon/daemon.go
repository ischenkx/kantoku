package daemon

import "context"

type Daemon struct {
	Settings Settings
	Func     func(ctx context.Context) error
	Type     string
}
