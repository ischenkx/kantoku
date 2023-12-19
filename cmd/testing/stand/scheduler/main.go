package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/processors/scheduler/dummy"
	"log/slog"
)

func main() {
	common.InitLogger()

	slog.Info("Starting...")
	err := dummy.NewProcessor(common.NewSystem(context.Background(), "scheduler-0")).Process(context.Background())
	if err != nil {
		slog.Error("failed", slog.String("error", err.Error()))
	}
}
