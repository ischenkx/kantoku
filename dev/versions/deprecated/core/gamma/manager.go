package gamma

import (
	"kantoku/common/data/bimap"
	"kantoku/common/data/deps"
	"kantoku/core/beta"
)

type Manager struct {
	betas        *beta.Manager
	dependencies deps.Deps
	beta2group   bimap.Bimap[string, string]
}
