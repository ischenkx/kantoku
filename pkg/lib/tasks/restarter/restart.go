package restarter

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
)

func Restart(ctx context.Context, system system.AbstractSystem, id string, infoCopiers ...InfoCopier) (newTaskID string, err error) {
	t, err := system.Task(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get task: %w", err)
	}

	rawStatus, ok := t.Info["status"]
	if !ok {
		return "", fmt.Errorf("no task status")
	}

	status, ok := rawStatus.(string)
	if !ok {
		return "", fmt.Errorf("invalid task status type (not string)")
	}

	if status != task.Statuses.Finished {
		return "", fmt.Errorf("task is not finished")
	}

	rawSubStatus, ok := t.Info["sub_status"]
	if !ok {
		return "", fmt.Errorf("no task sub_status")
	}

	subStatus, ok := rawSubStatus.(string)
	if !ok {
		return "", fmt.Errorf("invalid task sub_status type (not string)")
	}

	if subStatus != task.SubStatuses.Failed {
		return "", fmt.Errorf("task is not failed")
	}

	if rawRestarted, ok := t.Info["restarted"]; ok {
		if restarted, asserted := rawRestarted.(bool); restarted && asserted {
			return "", fmt.Errorf("already restarted")
		}
	}

	var restartRoot any = t.ID
	if parentRestartRoot, ok := t.Info["restart_root"]; ok {
		restartRoot = parentRestartRoot
	}

	modified, err := system.Tasks().UpdateWithProperties(ctx,
		map[string][]any{
			"info.restarted": {nil},
			"id":             {id},
		},
		map[string]any{
			"info.restarted":    true,
			"info.restart_root": restartRoot,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to update: %w", err)
	}

	if modified != 1 {
		return "", fmt.Errorf("no task updated")
	}

	newInfo := make(map[string]any)
	copyEssentialInfo(t.Info, newInfo)

	for _, copier := range infoCopiers {
		if copier == nil {
			continue
		}
		if err := copier(ctx, system, t, newInfo); err != nil {
			return "", fmt.Errorf("failed to copy info: %w", err)
		}
	}

	newTask, err := system.Spawn(ctx, task.New(
		task.WithInputs(t.Inputs...),
		task.WithOutputs(t.Outputs...),
		task.WithInfo(newInfo),
		task.WithProperty("restart_parent", t.ID),
		task.WithProperty("restart_root", restartRoot),
	))
	// TODO: rollback a restart (we might need another service that would persistently try to restart tasks)
	if err != nil {
		return "", fmt.Errorf("failed to spawn a new task: %w", err)
	}

	return newTask.ID, nil
}

type InfoCopier func(ctx context.Context, system system.AbstractSystem, oldTask task.Task, newTaskInfo map[string]any) error

func copyEssentialInfo(oldInfo, newInfo map[string]any) {
	typ, ok := oldInfo["type"]
	if ok {
		newInfo["type"] = typ
	}
	contextId, ok := oldInfo["context_id"]
	if ok {
		newInfo["context_id"] = contextId
	}
}
