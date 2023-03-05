package subtask

import (
	"context"
	"kantoku"
)

type Callback struct {
	subtask kantoku.Task
}

func (s Callback) Resolve(ctx context.Context, task kantoku.Task, app kantoku.Kantoku) (*kantoku.Argument, error) {
	//plugin, ok := app.PluginInstance(Id()).(Subtask)
	//if !ok {
	//	return nil, errors.New("could not cast plugin")
	//}
	taskId, err := app.New(ctx, s.subtask)
	if err != nil {
		return nil, err
	}
	task.Dependencies = append(task.Dependencies)
	return &kantoku.Argument{
		Type_: kantoku.ArgumentTypeTaskResult,
		Value: taskId,
	}, nil
}

func (s Callback) Register(ctx context.Context, task kantoku.Task, kantoku kantoku.Kantoku) error {
	return nil
}
