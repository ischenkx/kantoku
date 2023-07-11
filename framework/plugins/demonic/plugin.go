package demonic

import (
	"context"
	"kantoku/framework/infra/demon"
)

type Plugin struct {
	demons []demon.Demon
}

func New(demons ...demon.Demon) Plugin {
	return Plugin{demons: demons}
}

func (plugin Plugin) Demons(_ context.Context) []demon.Demon {
	return plugin.demons
}
