package main

import (
	"context"
	"fmt"
	"kantoku"
	"kantoku/unused/backend/stand/common"
	"time"
)

func main() {
	kan := kantoku.New(common.MakePlatform())
	result, err := kan.Spawn(
		context.Background(),
		kantoku.Task("factorial", []byte("12")),
	)

	if err != nil {
		fmt.Println("Failed to spawn task:", err)
		return
	}
	fmt.Println("Spawned task:", result.Task)

	time.Sleep(time.Second * 5)
	taskView := kan.Task(result.Task)
	data, err := taskView.Data(context.Background())
	if err != nil {
		fmt.Println("Failed to get task data:", err)
		return
	}
	fmt.Println("Data:", string(data))

	output, err := taskView.Result(context.Background())
	if err != nil {
		fmt.Println("Failed to get task result:", err)
		return
	}
	fmt.Println("Output:", string(output.Data))
}
