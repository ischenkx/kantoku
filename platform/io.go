package platform

import (
	"kantoku/common/data/kv"
	"kantoku/common/data/pool"
)

type Inputs[T Task] pool.Pool[T]
type Outputs kv.Database[string, Result]
