package scraper

import (
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn/future"
	"time"
)

type (
	DownloadPageInput struct {
		Url future.Future[string]
	}

	DownloadPageOutput struct {
		Page future.Future[[]byte]
	}

	DownloadPage struct {
		fn.Function[DownloadPage, DownloadPageInput, DownloadPageOutput]
	}
)

var (
	_ fn.AbstractFunction[DownloadPageInput, DownloadPageOutput] = (*DownloadPage)(nil)
)

func (f DownloadPage) Call(ctx *fn.Context, input DownloadPageInput) (output DownloadPageOutput, err error) {
	output.Page = future.FromValue([]byte("--- Sample Page ---"))

	time.Sleep(time.Second * 10)

	return output, nil
}
