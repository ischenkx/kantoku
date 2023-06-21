package postgres

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"kantoku/framework/plugins/depot/deps"
	"strings"
)

type Deps interface {
	deps.Deps
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
