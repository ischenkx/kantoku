package main

import (
	"context"
	http_tasks "github.com/ischenkx/kantoku/cmd/stand/sample_project/net/http"
	"github.com/ischenkx/kantoku/cmd/stand/sample_project/test"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http/oas"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn_d"
	"log"
)

func main() {
	rawClient, err := oas.NewClientWithResponses("http://localhost:8585")
	if err != nil {
		log.Fatal(err)
	}

	client := http.NewClient(rawClient)

	httpDoTask := http_tasks.Do{}
	httpDoTask.ID = "http.Do"
	httpDoSpec := fn_d.ToSpecification(httpDoTask.Function)

	randFailTask := test.RandFail{}
	randFailTask.ID = "test.RandFail"
	randFailSpec := fn_d.ToSpecification(randFailTask.Function)

	if err := client.Specifications().Add(context.Background(), httpDoSpec); err != nil {
		log.Fatal("failed to add a specification:", err)
	}

	if err := client.Specifications().Add(context.Background(), randFailSpec); err != nil {
		log.Fatal("failed to add a specification:", err)
	}
}
