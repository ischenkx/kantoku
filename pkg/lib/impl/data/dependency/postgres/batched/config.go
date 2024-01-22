package batched

import "time"

type Config struct {
	PollingInterval  time.Duration
	PollingBatchSize int
}
