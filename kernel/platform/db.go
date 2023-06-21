package platform

import "kantoku/common/data/kv"

type DB[T Task] kv.Database[string, T]
