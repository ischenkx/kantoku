package main

import (
	"context"
	"fmt"
	"kantoku"
	"kantoku/backend/stand/common"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func fact(x int) int {
	if x <= 1 {
		return 1
	}
	return x * fact(x-1)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	kan := common.MakeKantoku()

	x, err := kan.Futures().Make(context.Background(), "resource", nil)
	if err != nil {
		log.Fatal("failed to make a future:", err)
	}
	y, err := kan.Futures().Make(context.Background(), "resource", nil)
	if err != nil {
		log.Fatal("failed to make a future:", err)
	}
	z, err := kan.Futures().Make(context.Background(), "resource", nil)
	if err != nil {
		log.Fatal("failed to make a future:", err)
	}
	p, err := kan.Futures().Make(context.Background(), "resource", nil)
	if err != nil {
		log.Fatal("failed to make a future:", err)
	}
	q, err := kan.Futures().Make(context.Background(), "resource", nil)
	if err != nil {
		log.Fatal("failed to make a future:", err)
	}
	r, err := kan.Futures().Make(context.Background(), "resource", nil)
	if err != nil {
		log.Fatal("failed to make a future:", err)
	}

	mulTask, err := kan.Spawn(context.Background(),
		kantoku.Describe("mul").WithInputs(y.ID, z.ID).WithOutputs(p.ID),
	)
	if err != nil {
		log.Fatal("failed to spawn a task:", err)
	}

	factorialTask, err := kan.Spawn(context.Background(),
		kantoku.Describe("factorial").WithInputs(x.ID).WithOutputs(y.ID),
	)
	if err != nil {
		log.Fatal("failed to spawn a task:", err)
	}

	mulTask1, err := kan.Spawn(context.Background(),
		kantoku.Describe("mul").WithInputs(p.ID, q.ID).WithOutputs(r.ID),
	)
	if err != nil {
		log.Fatal("failed to spawn a task:", err)
	}

	fmt.Println("Factorial task:", factorialTask.Task)
	fmt.Println("Mul task:", mulTask.Task)
	fmt.Println("Mul task 1:", mulTask1.Task)
	fmt.Println("Resolving inputs")

	xVal := rand.Intn(20)
	zVal := rand.Intn(40)
	qVal := rand.Intn(40)

	fmt.Println("X:", xVal)
	fmt.Println("Z:", zVal)
	fmt.Println("EXPECTED:", fact(xVal)*zVal*qVal)

	if err := kan.Futures().Resolve(context.Background(), x.ID, []byte(strconv.Itoa(xVal))); err != nil {
		log.Fatal("failed to resolve x:", err)
	}
	if err := kan.Futures().Resolve(context.Background(), z.ID, []byte(strconv.Itoa(zVal))); err != nil {
		log.Fatal("failed to resolve z:", err)
	}
	if err := kan.Futures().Resolve(context.Background(), q.ID, []byte(strconv.Itoa(zVal))); err != nil {
		log.Fatal("failed to resolve q:", err)
	}
	for {
		time.Sleep(time.Second * 2)
		fmt.Println("Fetching results...")
		resolution, err := kan.Futures().Load(context.Background(), p.ID)
		if err != nil {
			fmt.Println("failed to load resolution:", err)
			continue
		}
		fmt.Println("Resolution:", string(resolution.Resource))
		break
	}
}
