package main

import (
	"context"
	"kantoku/backend/stand/common"
	"log"
)

func main() {
	kan, err := common.MakeKantoku()
	if err != nil {
		log.Fatal("failed to make kantoku:", err)
	}

	if err := kan.Infra().Deploy(context.Background()); err != nil {
		log.Fatal("failed to run demon:", err)
	}
}
