package main

import (
	"context"
	"github.com/samber/lo"
	"kantoku/core/l1"
	"kantoku/framework/executor/simple"
	"kantoku/testing/app/base"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kan, err := base.Generate(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	executor := simple.Executor{
		"reverse": func(data []byte) ([]byte, error) {
			return lo.Reverse(data), nil
		},
	}

	inputs := base.L1Inputs()
	outputs := base.L1Outputs(kan.Cells())

	if err := l1.New(inputs, outputs, executor, kan.Events()).Run(context.Background()); err != nil {
		log.Println("failed to run the l1 runner:", err)
	}
}
