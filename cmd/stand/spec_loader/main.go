package main

import (
	"context"
	http_tasks "github.com/ischenkx/kantoku/cmd/stand/sample_project/net/http"
	"github.com/ischenkx/kantoku/cmd/stand/sample_project/recursive"
	"github.com/ischenkx/kantoku/cmd/stand/sample_project/scraper"
	"github.com/ischenkx/kantoku/cmd/stand/sample_project/test"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/kantokuhttp"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/kantokuhttp/oas"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn"
	"log"
	"strings"
)

func main() {
	rawClient, err := oas.NewClientWithResponses("http://localhost:8585")
	if err != nil {
		log.Fatal(err)
	}

	client := kantokuhttp.NewClient(rawClient)

	register(client.Specifications(), http_tasks.Do{}.Function)
	register(client.Specifications(), test.RandFail{}.Function)
	register(client.Specifications(), recursive.A{}.Function)
	register(client.Specifications(), recursive.B{}.Function)
	register(client.Specifications(), scraper.Scrape{}.Function)
	register(client.Specifications(), scraper.DownloadPage{}.Function)
	register(client.Specifications(), scraper.ParsePage{}.Function)
	register(client.Specifications(), scraper.ExtractImages{}.Function)
}

func register[F fn.AbstractFunction[I, O], I, O any](storage *kantokuhttp.SpecificationStorage, function fn.Function[F, I, O]) {
	var f F
	e := fn.NewExecutor(f)

	spec := fn.ToSpecification(function)
	spec.ID = e.Type()
	spec.ID = strings.TrimPrefix(spec.ID, "github.com/ischenkx/kantoku/cmd/stand/")

	if err := storage.Add(context.Background(), spec); err != nil {
		log.Fatal("failed to add a specification:", err)
	}
}
