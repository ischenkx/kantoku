package platform

import (
	"kantoku/common/data/kv"
	"kantoku/common/data/pool"
)

type Inputs pool.Pool[string]
type Outputs kv.Database[string, Result]
