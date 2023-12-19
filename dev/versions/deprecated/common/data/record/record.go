package record

import (
	"github.com/samber/lo"
	"strings"
)

type Record map[string]any

func (r Record) String() string {
	entries := make([]string, 0, len(r))
	for name, value := range r {
		entries = append(entries, Entry{name, value}.String())
	}
	return strings.Join(entries, "; ")
}

func (r Record) Copy() Record {
	copied := R{}
	for key, value := range r {
		copied[key] = value
	}
	return copied
}

func (r Record) AsEntries() []Entry {
	return lo.MapToSlice(r, func(key string, value any) Entry {
		return Entry{key, value}
	})
}
