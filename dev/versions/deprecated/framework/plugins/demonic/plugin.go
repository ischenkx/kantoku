package demonic

import (
	"kantoku/framework/infra"
)

type Plugin struct {
	demons []infra.Demon
}

func New(demons ...infra.Demon) Plugin {
	return Plugin{demons: demons}
}

func (plugin Plugin) Demons() []infra.Demon {
	return plugin.demons
}
