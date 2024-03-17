package dumb

import (
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/samber/lo"
	"sort"
)

func matches(r record.R, filter record.R) bool {
	if r == nil {
		r = record.R{}
	}

	for r

	return true
}

func mask(r record.R, masks []record.Mask) record.R {
	if r == nil {
		r = record.R{}
	}

	includes := lo.Filter(masks, func(mask record.Mask, _ int) bool {
		return mask.Operation == record.IncludeMask
	})

	excludes := lo.Filter(masks, func(mask record.Mask, _ int) bool {
		return mask.Operation == record.ExcludeMask
	})

	if len(includes) > 0 && len(excludes) > 0 {
		return r
	}

	if len(excludes) > 0 {
		for _, exclude := range excludes {
			delete(r, exclude.PropertyPattern)
		}
	}

	if len(includes) > 0 {
		tmp := map[string]any{}
		for _, include := range includes {
			tmp[include.PropertyPattern] = r[include.PropertyPattern]
		}
		r = tmp
	}

	return r
}

func keyMask(r record.R, keys []string) []record.E {
	var res []record.E
	for _, key := range keys {
		res = append(res, record.E{key, getField(r, key)})
	}

	return res
}

func sorted(records []record.R, sorters []record.Sorter) {
	customSort := func(i, j int) bool {
		for _, sorter := range sorters {
			v1 := getField(records[i], sorter.Key)
			v2 := getField(records[j], sorter.Key)

			if sorter.Ordering == record.Asc {
				if le(v1, v2) {
					return true
				} else if ge(v1, v2) {
					return false
				}
			} else if sorter.Ordering == record.Desc {
				if le(v1, v2) {
					return false
				} else if ge(v1, v2) {
					return true
				}
			}
		}

		return i < j
	}

	// Sort the records using the custom sorting function
	sort.Slice(records, customSort)
}

func getField(record record.R, key string) any {
	if record == nil {
		return nil
	}

	return record[key]
}

func le(x, y any) bool {
	if x == nil && y != nil {
		return true
	}
	if y == nil && x != nil {
		return false
	}
	switch v := x.(type) {
	case int:
		if w, ok := y.(int); ok {
			return v < w
		}
	case int8:
		if w, ok := y.(int8); ok {
			return v < w
		}
	case int16:
		if w, ok := y.(int16); ok {
			return v < w
		}
	case int32:
		if w, ok := y.(int32); ok {
			return v < w
		}
	case int64:
		if w, ok := y.(int64); ok {
			return v < w
		}
	case uint:
		if w, ok := y.(uint); ok {
			return v < w
		}
	case uint8:
		if w, ok := y.(uint8); ok {
			return v < w
		}
	case uint16:
		if w, ok := y.(uint16); ok {
			return v < w
		}
	case uint32:
		if w, ok := y.(uint32); ok {
			return v < w
		}
	case uint64:
		if w, ok := y.(uint64); ok {
			return v < w
		}
	case uintptr:
		if w, ok := y.(uintptr); ok {
			return v < w
		}
	case float32:
		if w, ok := y.(float32); ok {
			return v < w
		}
	case float64:
		if w, ok := y.(float64); ok {
			return v < w
		}
	case string:
		if w, ok := y.(string); ok {
			return v < w
		}
	}
	return false
}

func ge(x, y any) bool {
	if x == nil && y != nil {
		return false
	}
	if y == nil && x != nil {
		return true
	}
	switch v := x.(type) {
	case int:
		if w, ok := y.(int); ok {
			return v > w
		}
	case int8:
		if w, ok := y.(int8); ok {
			return v > w
		}
	case int16:
		if w, ok := y.(int16); ok {
			return v > w
		}
	case int32:
		if w, ok := y.(int32); ok {
			return v > w
		}
	case int64:
		if w, ok := y.(int64); ok {
			return v > w
		}
	case uint:
		if w, ok := y.(uint); ok {
			return v > w
		}
	case uint8:
		if w, ok := y.(uint8); ok {
			return v > w
		}
	case uint16:
		if w, ok := y.(uint16); ok {
			return v > w
		}
	case uint32:
		if w, ok := y.(uint32); ok {
			return v > w
		}
	case uint64:
		if w, ok := y.(uint64); ok {
			return v > w
		}
	case uintptr:
		if w, ok := y.(uintptr); ok {
			return v > w
		}
	case float32:
		if w, ok := y.(float32); ok {
			return v > w
		}
	case float64:
		if w, ok := y.(float64); ok {
			return v > w
		}
	case string:
		if w, ok := y.(string); ok {
			return v > w
		}
	}
	return false
}

type Trie struct {
	value    any
	children []*Trie
}

func newTrie(value any) *Trie {
	return &Trie{
		value: value,
	}
}

func (trie *Trie) Insert(path []any) {
	if len(path) == 0 {
		return
	}
	var child *Trie
	for _, c := range trie.children {
		if c.value == path[0] {
			child = c
			break
		}
	}
	if child == nil {
		child = newTrie(path[0])
		trie.children = append(trie.children, child)
	}
	child.Insert(path[1:])
}
func (trie *Trie) Exists(path []any) bool {
	if len(path) == 0 {
		return true
	}
	for _, c := range trie.children {
		if c.value == path[0] {
			return c.Exists(path[1:])
		}
	}
	return false
}
