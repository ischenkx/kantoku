package batched

import (
	"kantoku/common/data/deps"
)

type DependencyModel struct {
	ID       string `bson:"id"`
	Resolved bool   `bson:"last_resolution"`
}

func (model DependencyModel) AsEntity() deps.Dependency {
	return deps.Dependency{
		ID:       model.ID,
		Resolved: model.Resolved,
	}
}

const (
	InitializingStatus = "initializing"
	WaitingStatus      = "waiting"
	SchedulingStatus   = "scheduling"
	ScheduledStatus    = "scheduled"
)

type GroupModel struct {
	ID      string `bson:"id"`
	Pending int    `bson:"dependencies"`
	Status  string `bson:"status"`
}
