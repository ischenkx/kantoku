package main

import (
	"context"
	"github.com/samber/lo"
	"kantoku"
	"kantoku/unused/backend/executor/evexec"
	"kantoku/unused/backend/stand/common"
	"log"
	"strconv"
)

type Runner struct{}

func (r *Runner) Run(ctx context.Context, task kantoku.TaskInstance) ([]byte, error) {
	log.Println("Running:", task.ID())
	switch task.Type {
	case "factorial":
		x, err := strconv.Atoi(string(task.Data))
		if err != nil {
			return nil, err
		}
		return []byte(strconv.Itoa(factorial(x))), nil
	case "reverse":
		return lo.Reverse(task.Data), nil
	}

	return nil, nil
}

func main() {
	executor := evexec.Builder[kantoku.TaskInstance]{
		Runner:   &Runner{},
		Platform: common.MakePlatform(),
		Resolver: evexec.ConstantResolver("TEST_STAND_EVENTS"),
	}.Build()

	log.Println("Starting...")
	if err := executor.Run(context.Background()); err != nil {
		log.Println("failed to run:", err)
	}
}

func factorial(x int) int {
	if x <= 1 {
		return 1
	}
	return x * factorial(x-1)
}
