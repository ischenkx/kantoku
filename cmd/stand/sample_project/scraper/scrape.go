package scraper

import (
	"fmt"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn/future"
)

type (
	ScrapeInput struct {
		Url future.Future[string]
	}

	ScrapeOutput struct {
		Result future.Future[map[string]any]
		Images future.Future[[]string]
	}

	Scrape struct {
		fn.Function[Scrape, ScrapeInput, ScrapeOutput]
	}
)

var (
	_ fn.AbstractFunction[ScrapeInput, ScrapeOutput] = (*Scrape)(nil)
)

func (f Scrape) Call(ctx *fn.Context, input ScrapeInput) (output ScrapeOutput, err error) {
	downloadPageResult, err := fn.Sched[DownloadPage](ctx, DownloadPageInput{
		Url: input.Url,
	})
	if err != nil {
		return output, fmt.Errorf("failed to download page: %w", err)
	}

	extractImagesResult, err := fn.Sched[ExtractImages](ctx, ExtractImagesInput{Page: downloadPageResult.Page})
	if err != nil {
		return output, fmt.Errorf("failed to extract images: %w", err)
	}

	parsePageResult, err := fn.Sched[ParsePage](ctx, ParsePageInput{Page: downloadPageResult.Page})
	if err != nil {
		return output, fmt.Errorf("failed to parse page: %w", err)
	}

	return ScrapeOutput{
		Result: parsePageResult.Result,
		Images: extractImagesResult.Images,
	}, nil
}
