package main

import (
	"fmt"
	"github.com/ischenkx/kantoku"
	webserver "github.com/ischenkx/kantoku/pkg/lib/connector/web/server"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
)

func main() {
	apiHost := "localhost"
	if host, ok := os.LookupEnv("API_HOST"); ok {
		apiHost = host
	}

	kanto, err := kantoku.Connect(fmt.Sprintf("http://%s:8080", apiHost))
	if err != nil {
		log.Fatal("failed to connect to kantoku:", err)
		return
	}

	server := webserver.Server{
		System: kanto,
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
