package demons

import (
	"context"
	"kantoku/framework/infra"
)

type List []infra.Demon

func (l List) Demons(_ context.Context) []infra.Demon {
	return l
}
