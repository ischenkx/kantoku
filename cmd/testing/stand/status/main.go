package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	codec "github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/processors/executor"
	"github.com/ischenkx/kantoku/pkg/processors/status"
	"log/slog"
)

func main() {
	common.InitLogger()

	slog.Info("Starting...")
	proc := status.NewProcessor(common.NewSystem(context.Background(), "status-0"), codec.JSON[executor.Result]())
	if err := proc.Process(context.Background()); err != nil {
		slog.Error("failed", slog.String("error", err.Error()))
	}
}
