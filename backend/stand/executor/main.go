package main

import (
	"context"
	"kantoku"
	evexec2 "kantoku/backend/executor/evexec"
	"kantoku/backend/stand/common"
	"kantoku/kernel"
	"log"
	"strconv"
	"time"
)

type Runner struct {
	kantoku *kantoku.Kantoku
}

func fact(x int) int {
	if x <= 1 {
		return 1
	}
	return x * fact(x-1)
}

func (r *Runner) Run(ctx context.Context, raw kernel.Task) ([]byte, error) {
	log.Println("running task:", raw.ID(), raw.Type)

	task := r.kantoku.Task(raw.ID())
	inputs, err := task.Inputs(ctx)
	if err != nil {
		log.Println("failed to load inputs:", err)
		return nil, err
	}

	outputs, err := task.Outputs(ctx)
	if err != nil {
		log.Println("failed to load outputs:", err)
		return nil, err
	}

	switch raw.Type {
	case "factorial":
		input, err := r.kantoku.Futures().Load(ctx, inputs[0])
		if err != nil {
			log.Println("failed to load input:", err)
			return nil, err
		}

		num, err := strconv.Atoi(string(input.Resource))
		if err != nil {
			log.Println("failed to cast input to a number:", err)
			return nil, err
		}

		res := fact(num)
		time.Sleep(time.Second * 10)
		err = r.kantoku.Futures().Resolve(ctx, outputs[0], []byte(strconv.Itoa(res)))
		if err != nil {
			log.Println("failed to resolve a given output future:", err)
			return nil, err
		}
	case "mul":
		input1, err := r.kantoku.Futures().Load(ctx, inputs[0])
		if err != nil {
			log.Println("failed to load input1:", err)
			return nil, err
		}

		input2, err := r.kantoku.Futures().Load(ctx, inputs[1])
		if err != nil {
			log.Println("failed to load input2:", err)
			return nil, err
		}

		num1, err := strconv.Atoi(string(input1.Resource))
		if err != nil {
			log.Println("failed to cast input1 to a number:", err)
			return nil, err
		}

		num2, err := strconv.Atoi(string(input2.Resource))
		if err != nil {
			log.Println("failed to cast input2 to a number:", err)
			return nil, err
		}

		res := num1 * num2
		err = r.kantoku.Futures().Resolve(ctx, outputs[0], []byte(strconv.Itoa(res)))
		if err != nil {
			log.Println("failed to resolve a given output future:", err)
			return nil, err
		}
	}

	return nil, nil
}

func main() {
	executor := evexec2.Builder[kernel.Task]{
		Runner:   &Runner{kantoku: common.MakeKantoku()},
		Platform: common.MakePlatform(),
		Resolver: evexec2.ConstantResolver("TEST_STAND_EVENTS"),
	}.Build()

	log.Println("Starting...")
	for i := 0; i < 12; i++ {
		go func() {
			if err := executor.Run(context.Background()); err != nil {
				log.Println("failed to run:", err)
			}
		}()
	}

	time.Sleep(time.Hour * 24)
}
