package main

import (
	"context"
	"flag"
	"fmt"
	"kantoku"
	"kantoku/backend/stand/common"
	"kantoku/common/data/future"
	"kantoku/common/util/futures"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func randomRequest(kan *kantoku.Kantoku) error {
	var x, y, z, p, q, r future.ID

	err := futures.Make(context.Background(), kan.Futures(), &x, &y, &z, &p, &q, &r)
	if err != nil {
		return fmt.Errorf("failed to make a future: %s", err)
	}

	_, err = kan.Tasks().Spawn(context.Background(),
		kantoku.Describe("mul").WithInputs(y, z).WithOutputs(p),
	)
	if err != nil {
		return fmt.Errorf("failed to spawn a task: %s", err)
	}

	_, err = kan.Tasks().Spawn(context.Background(),
		kantoku.Describe("factorial").WithInputs(x).WithOutputs(y),
	)
	if err != nil {
		return fmt.Errorf("failed to spawn a task: %s", err)
	}

	_, err = kan.Tasks().Spawn(context.Background(),
		kantoku.Describe("mul").WithInputs(p, q).WithOutputs(r),
	)
	if err != nil {
		return fmt.Errorf("failed to spawn a task: %s", err)
	}

	xVal := rand.Intn(20)
	zVal := rand.Intn(40)
	qVal := rand.Intn(40)

	if err := kan.Futures().OK(context.Background(), x, []byte(strconv.Itoa(xVal))); err != nil {
		return fmt.Errorf("failed to resolve x: %s", err)
	}
	if err := kan.Futures().OK(context.Background(), z, []byte(strconv.Itoa(zVal))); err != nil {
		return fmt.Errorf("failed to resolve z: %s", err)
	}
	if err := kan.Futures().OK(context.Background(), q, []byte(strconv.Itoa(qVal))); err != nil {
		return fmt.Errorf("failed to resolve q: %s", err)
	}
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		time.Sleep(time.Millisecond * 300)
		_, err := kan.Futures().Load(context.Background(), r)
		if err != nil {
			continue
		}
		break
	}
	return nil
}

var (
	TestsPtr = flag.Int("tests", 100, "amount of tests to run")
)

func main() {
	flag.Parse()

	TESTS := *TestsPtr

	log.Printf("Running %d tests\n", TESTS)
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

	for i := 0; i < TESTS; i++ {
		<-channel
		log.Printf("DONE %d / %d (avg: %s)", i+1, TESTS, time.Since(begin)/time.Duration(i+1))
	}

	fmt.Println("AVERAGE:", time.Since(begin)/time.Duration(TESTS))
}
