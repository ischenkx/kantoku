package manager

import "context"

type TaskToGroup interface {
	Save(ctx context.Context, task string, group string) error
	TaskByGroup(ctx context.Context, group string) (task string, err error)
	GroupByTask(ctx context.Context, task string) (group string, err error)
}
