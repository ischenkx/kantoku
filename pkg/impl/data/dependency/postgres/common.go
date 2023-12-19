package postgres

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/deps"
	"github.com/samber/lo"
	"strings"
)

type Deps interface {
	deps.Manager
	Run(ctx context.Context)
	InitTables(ctx context.Context) error
	DropTables(ctx context.Context) error
}

func FormatValues(ids ...string) string {
	values := strings.Join(
		lo.Map(ids, func(item string, _ int) string {
			return fmt.Sprintf("'%s'", item)
		}),
		", ")
	values = fmt.Sprintf("%s", values)
	return values
}
