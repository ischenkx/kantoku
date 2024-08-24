package main

import (
	"context"
	http_tasks "github.com/ischenkx/kantoku/cmd/stand/sample_project/net/http"
	"github.com/ischenkx/kantoku/cmd/stand/sample_project/recursive"
	"github.com/ischenkx/kantoku/cmd/stand/sample_project/scraper"
	"github.com/ischenkx/kantoku/cmd/stand/sample_project/test"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http/oas"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn"
	"log"
	"strings"
)

func main() {
	rawClient, err := oas.NewClientWithResponses("http://localhost:8585")
	if err != nil {
		log.Fatal(err)
	}

	client := http.NewClient(rawClient)

	register(client.Specifications(), http_tasks.Do{}.Function)
	register(client.Specifications(), test.RandFail{}.Function)
	register(client.Specifications(), recursive.A{}.Function)
	register(client.Specifications(), recursive.B{}.Function)
	register(client.Specifications(), scraper.Scrape{}.Function)
	register(client.Specifications(), scraper.DownloadPage{}.Function)
	register(client.Specifications(), scraper.ParsePage{}.Function)
	register(client.Specifications(), scraper.ExtractImages{}.Function)

	//fn.NewExecutor()
	//
	//httpDoTask := http_tasks.Do{}
	//httpDoTask.ID = httpDoTask.
	//httpDoSpec := fn.ToSpecification(httpDoTask.Function)
	//
	//randFailTask := test.RandFail{}
	//randFailTask.ID = "test.RandFail"
	//randFailSpec := fn.ToSpecification(randFailTask.Function)
	//
	//if err := client.Specifications().Add(context.Background(), httpDoSpec); err != nil {
	//	log.Fatal("failed to add a specification:", err)
	//}
	//
	//if err := client.Specifications().Add(context.Background(), randFailSpec); err != nil {
	//	log.Fatal("failed to add a specification:", err)
	//}
}

func register[F fn.AbstractFunction[I, O], I, O any](storage *http.SpecificationStorage, function fn.Function[F, I, O]) {
	var f F
	e := fn.NewExecutor(f)

	spec := fn.ToSpecification(function)
	spec.ID = e.Type()
	spec.ID = strings.TrimPrefix(spec.ID, "github.com/ischenkx/kantoku/cmd/stand/")

	if err := storage.Add(context.Background(), spec); err != nil {
		log.Fatal("failed to add a specification:", err)
	}
}
