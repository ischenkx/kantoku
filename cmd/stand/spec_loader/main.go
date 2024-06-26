package main

import (
	"context"
	http_tasks "github.com/ischenkx/kantoku/local/deprecated/sample_project/net/http"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http/oas"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/functional"
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

	httpDoSpec := functional.ToSpecification(httpDoTask.Task)

	if err := client.Specifications().Add(context.Background(), httpDoSpec); err != nil {
		log.Fatal("failed to add a specification:", err)
	}
}
