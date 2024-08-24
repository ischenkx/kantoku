package test

import (
	"fmt"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn/future"
	"math/rand"
)

type (
	RandFailInput struct{}

	RandFailOutput struct {
		Code future.Future[int]
	}

	RandFail struct {
		fn.Function[*RandFail, RandFailInput, RandFailOutput]
	}
)

var (
	_ fn.AbstractFunction[RandFailInput, RandFailOutput] = (*RandFail)(nil)
)

func (task *RandFail) Call(ctx *fn.Context, input RandFailInput) (output RandFailOutput, err error) {
	num := rand.Intn(100)
	if num < 50 {
		return RandFailOutput{}, fmt.Errorf("you lost: %d", num)
	}

	return RandFailOutput{Code: future.FromValue(num)}, nil
}
