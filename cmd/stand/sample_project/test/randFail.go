package test

import (
	"fmt"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn_d"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn_d/future"
	"math/rand"
)

type (
	RandFailInput struct{}

	RandFailOutput struct {
		Code future.Future[int]
	}

	RandFail struct {
		fn_d.Function[RandFail, RandFailInput, RandFailOutput]
	}
)

var (
	_ fn_d.AbstractFunction[RandFailInput, RandFailOutput] = (*RandFail)(nil)
)

func (task RandFail) Call(ctx *fn_d.Context, input RandFailInput) (output RandFailOutput, err error) {
	num := rand.Intn(100)
	if num < 50 {
		return RandFailOutput{}, fmt.Errorf("you lost: %d", num)
	}

	return RandFailOutput{Code: future.FromValue(num)}, nil
}
