package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/samber/lo"
	"io"
	"log/slog"
	"os"
	"time"
)

var Interval = time.Millisecond * 1000

func main() {
	common.InitLogger()
	ctx := context.Background()
	pech := common.NewSystem(ctx, "foreman-0")

	urlCorpus, err := loadUrls()
	if err != nil {
		slog.Error("failed to load urls",
			slog.String("error", err.Error()))
		return
	}

	ticker := time.NewTicker(Interval)
	defer ticker.Stop()

	for range ticker.C {
		url := lo.Sample[string](urlCorpus)

		slog.Info("registering an http request",
			slog.String("url", url))

		registerHttpRequest(ctx, url, pech)
	}
}

func registerHttpRequest(ctx context.Context, url string, sys *system.System) {

	resources, err := sys.
		Resources().
		Alloc(ctx, 2)
	if err != nil {
		slog.Error("failed to allocate resources",
			slog.String("error", err.Error()))
		return
	}

	urlResource, bodyResource := resources[0], resources[1]

	err = sys.
		Resources().
		Init(ctx,
			[]resource.Resource{
				{
					ID:   urlResource,
					Data: []byte(url),
				},
			})
	if err != nil {
		slog.Error("failed to initialize resources",
			slog.String("error", err.Error()))
		return
	}

	scheduledTask, err := sys.Spawn(context.Background(),
		task.Task{
			Inputs:  []resource.ID{urlResource},
			Outputs: []resource.ID{bodyResource},
			Info: record.R{
				"type": "http_request",
			},
		},
	)

	if err != nil {
		slog.Error("failed to register a new task",
			slog.String("error", err.Error()))
		return
	}

	slog.Info("registered a new task:",
		slog.String("id", scheduledTask.ID),
		slog.Any("inputs", scheduledTask.Inputs),
		slog.Any("outputs", scheduledTask.Outputs),
		slog.Any("info", scheduledTask.Info),
	)
}

func loadUrls() ([]string, error) {
	data, err := os.ReadFile("urls")
	if err != nil {
		return nil, fmt.Errorf("failed to read the file: %w", err)
	}

	reader := bufio.NewReader(bytes.NewReader(data))

	var urls []string

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("failed to read line: %w", err)
		}

		urls = append(urls, string(line))
	}

	return urls, nil
}
