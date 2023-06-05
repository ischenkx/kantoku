package main

import (
	"context"
	"kantoku/unused/backend/stand/common"
)

func main() {
	common.MakeDeps().Run(context.Background())
	<-context.Background().Done()
}
