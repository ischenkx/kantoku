package main

import (
	"fmt"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http/oas"
	webserver "github.com/ischenkx/kantoku/pkg/lib/gateway/web/server"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
)

func main() {
	apiHost := "localhost"
	if host, ok := os.LookupEnv("API_HOST"); ok {
		apiHost = host
	}

	rawClient, err := oas.NewClientWithResponses(fmt.Sprintf("http://%s:8080", apiHost))
	if err != nil {
		log.Fatal("failed to create a raw client:", err)
		return
	}

	server := webserver.Server{
		System: http.New(rawClient),
	}

	e := server.Echo()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	// Log all routes
	for _, route := range e.Routes() {
		fmt.Printf("%s %s\n", route.Method, route.Path)
	}

	if err := e.Start(":3030"); err != nil {
		log.Fatal("failed to start the server:", err)
	}
}
