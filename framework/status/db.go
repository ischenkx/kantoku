package status

import "context"

type DB interface {
	UpdateStatus(ctx context.Context, id string, status Status) error
}
