package roller

import "context"

type Saga struct {
	name    string
	actions []Action
}

func (saga *Saga) Actions() []Action {
	return saga.actions
}

func (saga *Saga) Name() string {
	return saga.name
}

func (saga *Saga) Do(action Action) *Saga {
	saga.actions = append(saga.actions, action)
	return saga
}

func (saga *Saga) Run(ctx context.Context) error {
	for _, action := range saga.Actions() {
		if action.Func != nil {
			action.Func
		}
	}
}
