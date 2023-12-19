package deprecated

import (
	"kantoku/common/data/future"
)

type Parametrization struct {
	Inputs  []future.ID
	Outputs []future.ID
	Static  []byte
}
