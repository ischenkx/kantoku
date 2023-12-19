package deprecated

import (
	"kantoku/framework/plugins/depot"
	"kantoku/framework/plugins/futdep"
	"kantoku/framework/plugins/taskdep"
)

type DependencyManager struct {
	taskdep *taskdep.Manager
	futdep  *futdep.Manager
	depot   *depot.Depot
}

func (manager DependencyManager) Futures() *futdep.Manager {
	return manager.futdep
}

func (manager DependencyManager) Tasks() *taskdep.Manager {
	return manager.taskdep
}

func (manager DependencyManager) Depot() *depot.Depot {
	return manager.depot
}
