package postgres

import "kantoku/common/deps"

type Group struct {
	Spec   deps.Group
	Status string
}
