package main

import (
	"context"
	"fmt"
	"kantoku"
	"kantoku/unused/backend/stand/common"
)

func main() {
	result, err := kantoku.New(common.MakePlatform()).Spawn(
		context.Background(),
		kantoku.Task("factorial", []byte("12")),
	)

	if err != nil {
		fmt.Println("Failed:", err)
		return
	}

	fmt.Println("Done:", result.Task)
}
