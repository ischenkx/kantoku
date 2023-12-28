package main

import (
	"cmp"
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/extensions/web/client"
	"github.com/ischenkx/kantoku/pkg/extensions/web/oas"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
	"github.com/samber/lo"
	"log"
	"slices"
	"time"
)

type statusCounter struct {
	status  string
	counter int
}

func main() {
	rawClient, err := oas.NewClientWithResponses("http://localhost:8080")
	if err != nil {
		log.Fatal("failed to create a raw client:", err)
		return
	}

	c := client.New(rawClient)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	fmt.Println("Starting!")

	for range ticker.C {
		fmt.Println("-------///-------")
		counters, err := countStatuses(c.Info(),
			task.InitializedStatus,
			task.ReadyStatus,
			task.CancelledStatus,
			task.ReceivedStatus,
			task.OkStatus,
			task.FailedStatus)
		if err != nil {
			fmt.Println("failed to count statuses:", err)
			continue
		}

		statusCounters := lo.Map(
			lo.Entries(counters),
			func(entry lo.Entry[string, int], _ int) statusCounter {
				return statusCounter{
					status:  entry.Key,
					counter: entry.Value,
				}
			})

		slices.SortFunc(
			statusCounters,
			func(a, b statusCounter) int {
				return cmp.Compare(a.status, b.status)
			})

		for _, sc := range statusCounters {
			fmt.Printf("# '%s' -> %d\n", sc.status, sc.counter)
		}
	}
}

func countStatuses(storage record.Set, statuses ...string) (map[string]int, error) {
	counters := map[string]int{}

	for _, status := range statuses {
		counter, err := storage.
			Filter(record.R{"status": status}).
			Cursor().
			Count(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to count status '%s': %w", status, err)
		}
		counters[status] = counter
	}

	return counters, nil
}
