package model

import "github.com/ischenkx/kantoku/pkg/common/dependency"

type DependencyModel struct {
	ID     string
	Status string
}

func (model DependencyModel) AsEntity() dependency.Dependency {
	return dependency.Dependency{
		ID:     model.ID,
		Status: dependency.Status(model.Status),
	}
}

type GroupModel struct {
	ID      string `bson:"id"`
	Pending int    `bson:"dependencies"`
	Status  string `bson:"status"`
}
