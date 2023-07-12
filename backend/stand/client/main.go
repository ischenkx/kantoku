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

func randomRequest(kan *kantoku.Kantoku) error {
	x, err := kan.Futures().Make(context.Background(), "resource", nil)
	if err != nil {
		return fmt.Errorf("failed to make a future: %s", err)
	}
	y, err := kan.Futures().Make(context.Background(), "resource", nil)
	if err != nil {
		return fmt.Errorf("failed to make a future: %s", err)
	}
	z, err := kan.Futures().Make(context.Background(), "resource", nil)
	if err != nil {
		return fmt.Errorf("failed to make a future: %s", err)
	}
	p, err := kan.Futures().Make(context.Background(), "resource", nil)
	if err != nil {
		return fmt.Errorf("failed to make a future: %s", err)
	}
	q, err := kan.Futures().Make(context.Background(), "resource", nil)
	if err != nil {
		return fmt.Errorf("failed to make a future: %s", err)
	}
	r, err := kan.Futures().Make(context.Background(), "resource", nil)
	if err != nil {
		return fmt.Errorf("failed to make a future: %s", err)
	}

	_, err = kan.Spawn(context.Background(),
		kantoku.Describe("mul").WithInputs(y.ID, z.ID).WithOutputs(p.ID),
	)
	if err != nil {
		return fmt.Errorf("failed to spawn a task: %s", err)
	}

	_, err = kan.Spawn(context.Background(),
		kantoku.Describe("factorial").WithInputs(x.ID).WithOutputs(y.ID),
	)
	if err != nil {
		return fmt.Errorf("failed to spawn a task: %s", err)
	}

	_, err = kan.Spawn(context.Background(),
		kantoku.Describe("mul").WithInputs(p.ID, q.ID).WithOutputs(r.ID),
	)
	if err != nil {
		return fmt.Errorf("failed to spawn a task: %s", err)
	}

	xVal := rand.Intn(20)
	zVal := rand.Intn(40)
	qVal := rand.Intn(40)

	if err := kan.Futures().Resolve(context.Background(), x.ID, []byte(strconv.Itoa(xVal))); err != nil {
		return fmt.Errorf("failed to resolve x: %s", err)
	}
	if err := kan.Futures().Resolve(context.Background(), z.ID, []byte(strconv.Itoa(zVal))); err != nil {
		return fmt.Errorf("failed to resolve z: %s", err)
	}
	if err := kan.Futures().Resolve(context.Background(), q.ID, []byte(strconv.Itoa(qVal))); err != nil {
		return fmt.Errorf("failed to resolve q: %s", err)
	}
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		time.Sleep(time.Millisecond * 300)
		_, err := kan.Futures().Load(context.Background(), r.ID)
		if err != nil {
			continue
		}
		break
	}
	return nil
}

const TESTS = 100

func main() {
	begin := time.Now()
	defer func() {
		fmt.Println("ELAPSED:", time.Since(begin))
	}()
	rand.Seed(time.Now().UnixNano())
	kan, err := common.MakeKantoku()
	if err != nil {
		log.Fatal(err)
	}

	channel := make(chan int64, TESTS)
	for i := 0; i < TESTS; i++ {
		go func() {
			begin := time.Now()
			if err := randomRequest(kan); err != nil {
				log.Println("Failure!")
			}
			channel <- int64(time.Since(begin))
		}()
	}

	total := 0.0

	for i := 0; i < TESTS; i++ {
		res := <-channel
		total += float64(res)
		log.Printf("DONE %d / %d (avg: %s, current: %s)", i+1, TESTS, time.Duration(total/float64(i+1)), time.Duration(res))
	}

	fmt.Println("AVERAGE:", time.Duration(int64(total/TESTS)))
}
