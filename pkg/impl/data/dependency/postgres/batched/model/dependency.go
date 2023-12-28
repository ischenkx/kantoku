package model

import "github.com/ischenkx/kantoku/pkg/common/data/deps"

type DependencyModel struct {
	ID     string
	Status string
}

func (model DependencyModel) AsEntity() deps.Dependency {
	return deps.Dependency{
		ID:     model.ID,
		Status: deps.Status(model.Status),
	}
}

type GroupModel struct {
	ID      string `bson:"id"`
	Pending int    `bson:"dependencies"`
	Status  string `bson:"status"`
}
