package scraper

import (
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn/future"
)

type (
	ExtractImagesInput struct {
		Page future.Future[[]byte]
	}

	ExtractImagesOutput struct {
		Images future.Future[[]string]
	}

	ExtractImages struct {
		fn.Function[ExtractImages, ExtractImagesInput, ExtractImagesOutput]
	}
)

var (
	_ fn.AbstractFunction[ExtractImagesInput, ExtractImagesOutput] = (*ExtractImages)(nil)
)

func (f ExtractImages) Call(ctx *fn.Context, input ExtractImagesInput) (output ExtractImagesOutput, err error) {
	output.Images = future.FromValue([]string{
		"url1",
		"url2",
	})

	return output, nil
}
