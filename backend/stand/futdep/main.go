package main

import (
	"context"
	"kantoku/backend/stand/common"
	"kantoku/framework/plugins/futdep"
	"log"
)

func main() {
	manager := futdep.NewManager(
		common.MakeDeps(),
		common.MakeFutDepDB(),
	)
	queue := common.MakeFutureResolutionQueue()
	channel, err := queue.Read(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for fut := range channel {
		if err := manager.ResolveFuture(context.Background(), fut); err != nil {
			log.Println("failed to resolve future:", err)
		}
	}
}
