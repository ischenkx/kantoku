package main

import (
	"context"
	"kantoku"
	"kantoku/core/l2"
	"kantoku/impl/common/pool/func"
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

	inputs := base.L1Inputs()

	runner := l2.New[kantoku.Task](
		kan.Tasks(),
		kan.Events(),
		inputs,
	)

	err = kan.Depot().Process(
		ctx,
		funcpool.NewWriter[string](runner.Run),
	)
	if err != nil {
		log.Println("failed to run the depot processor:", err)
		return
	}
}
