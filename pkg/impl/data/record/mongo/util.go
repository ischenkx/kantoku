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

func makeRecordFilter(rec record.R) (filter bson.M, err error) {
	for key, value := range rec {
		filter, err := makeValueFilter(value)
		if err != nil {
			return nil, fmt.Errorf("failed to make a filter for '%s': %w", key, err)
		}
		rec[key] = filter
	}
	return
}

func makeValueFilter(value any) (any, error) {
	if operation, ok := value.(ops.Operation); ok {
		return operation2filter(operation)
	}

	return bson.M{"$eq": value}, nil
}

func operation2filter(op ops.Operation) (any, error) {
	switch op.Type {
	case ops.InOp:
		list, ok := op.Data.([]any)
		if !ok {
			return nil, fmt.Errorf("'in' operation expects a list as an argument")
		}

		return bson.M{"$in": list}, nil
	case ops.LessOp:
		return bson.M{"$lt": op.Data}, nil
	case ops.LessOrEqOp:
		return bson.M{"$lte": op.Data}, nil
	case ops.GreaterOp:
		return bson.M{"$gt": op.Data}, nil
	case ops.GreaterOrEq:
		return bson.M{"$gte": op.Data}, nil
	case ops.LikeOp:
		return nil, fmt.Errorf("like operation is not supported")
	case ops.ContainsOp:
		list, ok := op.Data.([]any)
		if !ok {
			return nil, fmt.Errorf("'contains' operation expects a list as an argument")
		}

		subFilters := make([]any, 0, len(list))
		for _, item := range list {
			subFilter, err := makeValueFilter(item)
			if err != nil {
				return nil, fmt.Errorf("failed to make a sub-filter for '%s': %w", item, err)
			}
			subFilters = append(subFilters, subFilter)
		}

		return bson.M{"$all": subFilters}, nil
	case ops.NotOp:
		subFilter, err := makeValueFilter(op.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to make a sub-filter for '%s': %w", op.Data, err)
		}
		return bson.M{"$not": subFilter}, nil
	case ops.AndOp:
		list, ok := op.Data.([]any)
		if !ok {
			return nil, fmt.Errorf("'and' operation expects a list as an argument")
		}

		subFilters := make([]any, 0, len(list))
		for _, item := range list {
			subFilter, err := makeValueFilter(item)
			if err != nil {
				return nil, fmt.Errorf("failed to make a sub-filter for '%s': %w", item, err)
			}
			subFilters = append(subFilters, subFilter)
		}

		return bson.M{"$and": subFilters}, nil
	case ops.OrOp:
		list, ok := op.Data.([]any)
		if !ok {
			return nil, fmt.Errorf("'or' operation expects a list as an argument")
		}

		subFilters := make([]any, 0, len(list))
		for _, item := range list {
			subFilter, err := makeValueFilter(item)
			if err != nil {
				return nil, fmt.Errorf("failed to make a sub-filter for '%s': %w", item, err)
			}
			subFilters = append(subFilters, subFilter)
		}

		return bson.M{"$or": subFilters}, nil
	default:
		return nil, fmt.Errorf("unknown operation '%s'", op.Type)
	}
}
