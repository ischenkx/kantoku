package demons

import (
	"context"
	"kantoku/framework/infra/demon"
)

type List []demon.Demon

func (l List) Demons(_ context.Context) []demon.Demon {
	return l
}
