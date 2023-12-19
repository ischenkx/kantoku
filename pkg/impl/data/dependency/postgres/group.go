package postgres

import "github.com/ischenkx/kantoku/pkg/common/data/deps"

type Group struct {
	Spec   deps.Group
	Status string
}
