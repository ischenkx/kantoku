package mongorec

import (
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/common/data/record/ops"
	"go.mongodb.org/mongo-driver/bson"
)

func record2bson(r record.R) bson.M {
	return bson.M(r)
}

func bson2record(m bson.M) record.R {
	return record.R(m)
}

func unwrapUpdateRecord(rec record.R) record.R {
	collector := record.R{}
	_unwrapUpdateRecord("", rec, collector)

	return collector
}

func _unwrapUpdateRecord(path string, value any, collector record.R) {
	switch rec := value.(type) {
	case record.R:
		for key, val := range rec {
			newPath := key
			if len(path) > 0 {
				newPath = fmt.Sprintf("%s.%s", path, key)
			}
			_unwrapUpdateRecord(newPath, val, collector)
		}
	case map[string]any:
		_unwrapUpdateRecord(path, record.R(rec), collector)
	default:
		if len(path) > 0 {
			collector[path] = value
		}
	}
}

func makeRecordFilter(rec record.R) (filter bson.M, err error) {
	filter = bson.M{}
	for key, value := range rec {
		err := applyToFilter(filter, key, value)
		if err != nil {
			return nil, fmt.Errorf("failed to make a filter for '%s': %w", key, err)
		}
	}
	return
}

func applyToFilter(filter bson.M, key string, value any) error {
	operation, ok := value.(ops.Operation)
	if !ok {
		operation = ops.Eq(value)
	}
	return applyOperation(filter, key, operation)
}

func applyOperation(filter bson.M, key string, op ops.Operation) error {
	switch op.Type {
	case ops.InOp:
		list, ok := op.Data.([]any)
		if !ok {
			return fmt.Errorf("'in' operation expects a list as an argument")
		}

		filter[key] = bson.M{"$in": list}
	case ops.LessOp:
		filter[key] = bson.M{"$lt": op.Data}
	case ops.LessOrEqOp:
		filter[key] = bson.M{"$lte": op.Data}
	case ops.GreaterOp:
		filter[key] = bson.M{"$gt": op.Data}
	case ops.GreaterOrEq:
		filter[key] = bson.M{"$gte": op.Data}
	case ops.LikeOp:
		return fmt.Errorf("like operation is not supported")
	case ops.ContainsOp:
		list, ok := op.Data.([]any)
		if !ok {
			return fmt.Errorf("'contains' operation expects a list as an argument")
		}

		filter[key] = bson.M{"$all": list}
	case ops.NotOp:
		subFilter := bson.M{}
		err := applyToFilter(subFilter, "$not", op.Data)
		if err != nil {
			return fmt.Errorf("failed to make a sub-filter for '%s': %w", op.Data, err)
		}
		filter[key] = subFilter
	case ops.AndOp:
		list, ok := op.Data.([]any)
		if !ok {
			return fmt.Errorf("'and' operation expects a list as an argument")
		}

		for _, item := range list {
			subFilter := bson.M{}

			err := applyToFilter(subFilter, key, item)
			if err != nil {
				return fmt.Errorf("failed to make a sub-filter for '%s': %w", item, err)
			}

			and, ok := filter["$and"]
			if !ok {
				and = []any{}
				filter["$and"] = and
			}

			typedAnd := and.([]any)
			typedAnd = append(typedAnd, subFilter)
			filter["$and"] = typedAnd
		}
	case ops.OrOp:
		list, ok := op.Data.([]any)
		if !ok {
			return fmt.Errorf("'or' operation expects a list as an argument")
		}

		for _, item := range list {
			subFilter := bson.M{}

			err := applyToFilter(subFilter, key, item)
			if err != nil {
				return fmt.Errorf("failed to make a sub-filter for '%s': %w", item, err)
			}

			or, ok := filter["$or"]
			if !ok {
				or = []any{}
				filter["$or"] = or
			}

			typedOr := or.([]any)
			typedOr = append(typedOr, subFilter)
			filter["$or"] = typedOr
		}
	case ops.EqOp:
		filter[key] = bson.M{"$eq": op.Data}
	default:
		return fmt.Errorf("unknown operation '%s'", op.Type)
	}

	return nil
}
