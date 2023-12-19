package postgres

import (
	"kantoku/common/data/deps"
)

type Group struct {
	Spec   deps.Group
	Status string
}
