package postgredeps

import (
	"kantoku/unused/backend/framework/depot/deps"
)

type Group struct {
	Spec   deps.Group
	Status string
}
