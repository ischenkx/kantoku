package main

import (
	"context"
	"kantoku/unused/backend/stand/common"
	"log"
)

func main() {
	if err := common.MakeInputs().Process(context.Background()); err != nil {
		log.Println("failed to process:", err)
	}
}
