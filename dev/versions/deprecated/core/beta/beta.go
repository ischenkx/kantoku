package beta

import (
	"context"
	"kantoku/core/alpha"
)

type Beta struct {
	manager *Manager
	id      string
}

func (beta Beta) ID() string {
	return beta.id
}

func (beta Beta) Info(ctx context.Context) (Info, error) {

}

func (beta Beta) Alpha(ctx context.Context) (alpha.Alpha, error) {

}
