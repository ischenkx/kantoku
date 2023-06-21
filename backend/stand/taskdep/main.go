package main

import (
	"context"
	"kantoku/backend/stand/common"
	"kantoku/framework/plugins/taskdep"
	"log"
)

func main() {

	updater := taskdep.NewUpdater(
		common.MakeBroker(),
		taskdep.NewManager(common.MakeDeps(), common.MakeTaskDepDB()),
		"TEST_STAND_EVENTS",
	)
	if err := updater.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
