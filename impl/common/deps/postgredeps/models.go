package postgredeps

import "kantoku/common/deps"

type DependencyModel struct {
	ID             string `bson:"id"`
	LastResolution int64  `bson:"last_resolution"`
}

func (model DependencyModel) AsEntity() deps.Dependency {
	return deps.Dependency{
		ID:             model.ID,
		LastResolution: model.LastResolution,
	}
}

const (
	WaitingStatus    = "waiting"
	SchedulingStatus = "scheduling"
	ScheduledStatus  = "scheduled"
)

type GroupModel struct {
	ID           string          `bson:"id"`
	Dependencies map[string]bool `bson:"dependencies"`
	Status       string
	Hash         string
}

func (model GroupModel) AsEntity() deps.Group {
	return deps.Group{
		ID:           model.ID,
		Dependencies: model.Dependencies,
	}
}
