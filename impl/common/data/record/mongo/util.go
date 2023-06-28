package mongorec

import (
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"kantoku/common/data/record"
)

func bsonifyEntries(entries []record.Entry) bson.M {
	return lo.SliceToMap[record.Entry, string, any](entries, func(entry record.Entry) (string, any) {
		return entry.Name, entry.Value
	})
}

func debsonifyEntries(m bson.M) []record.Entry {
	return lo.MapToSlice[string, any, record.Entry](m, func(name string, value any) record.Entry {
		return record.Entry{
			Name:  name,
			Value: value,
		}
	})
}

func bsonifyRecord(r record.R) bson.M {
	return bson.M(r)
}

func debsonifyRecord(m bson.M) record.R {
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
