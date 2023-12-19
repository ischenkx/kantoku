package alpha

import (
	"kantoku/common/data/kv"
	"kantoku/common/data/pool"
)

type Pool pool.Pool[string]
type ResultStorage kv.Database[string, Result]
type Storage kv.Database[string, []byte]
