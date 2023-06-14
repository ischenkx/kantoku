package kantoku

import "kantoku/framework/future"

type Parametrization struct {
	Inputs  []future.ID
	Outputs []future.ID
	Static  []byte
}
