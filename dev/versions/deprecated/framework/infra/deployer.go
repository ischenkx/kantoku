package infra

import (
	"context"
)

type Deployer interface {
	Deploy(ctx context.Context, demons ...Demon) error
}
