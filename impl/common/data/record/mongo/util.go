package mongorec

import (
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"kantoku/common/data/record"
)

func record2bson(r record.R) bson.M {
	return bson.M(r)
}

func bson2record(m bson.M) record.R {
	return record.R(m)
}

func makeFilter(filters [][]record.E) bson.M {
	conj := lo.FilterMap(filters, func(disj []record.E, _ int) (bson.M, bool) {
		if len(disj) == 0 {
			return nil, false
		}

		return bson.M{
			"$or": lo.Map(disj, func(e record.E, _ int) bson.M {
				return bson.M{e.Name: e.Value}
			}),
		}, true
	})

	if len(conj) == 0 {
		return bson.M{}
	}

	return bson.M{"$and": conj}
}
