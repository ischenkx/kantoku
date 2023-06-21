package postgres

import (
	"kantoku/framework/plugins/depot/deps"
)

type Group struct {
	Spec   deps.Group
	Status string
}
