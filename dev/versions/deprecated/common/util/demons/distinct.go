package demons

import (
	"kantoku/framework/infra"
)

func Distinct(demons ...infra.Demon) (result []infra.Demon) {
filter:
	for _, demon := range demons {
		for _, existing := range result {
			if demon.Eq(existing) {
				continue filter
			}
		}
		result = append(result, demon)
	}
	return
}
