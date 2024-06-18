package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ischenkx/kantoku/cmd/stand/utils"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/lib/platform"
	"os"
)

var resourceId = flag.String("resource", "", "Resource ID")
var filePath = flag.String("file", "", "File path")
var useJson = flag.Bool("json", false, "Parse as json")

func main() {
	flag.Parse()

	if *resourceId == "" {
		fmt.Println("You must specify a resource ID")
		return
	}

	if *filePath == "" {
		fmt.Println("You must specify a file path")
		return
	}

	ctx := context.Background()
	cfg := utils.LoadConfig()
	logger := utils.GetLogger(os.Stdout, "load_resource")

	sys, err := platform.BuildSystem(ctx, logger, cfg.Core.System)
	if err != nil {
		fmt.Println("failed to build system:", err)
		return
	}

	resources, err := sys.Resources().Load(ctx, *resourceId)
	if err != nil {
		fmt.Println("failed to load resources:", err)
		return
	}

	res := resources[0]

	if res.Status != resource.Ready {
		fmt.Printf("Resource is not ready (status=%s)\n", res.Status)
		return
	}

	data := res.Data

	if *useJson {
		var parsed string
		if err := json.Unmarshal(data, &parsed); err != nil {
			fmt.Println("failed to parse json:", err)
			return
		}
		data = []byte(parsed)
	}

	if err := os.WriteFile(*filePath, data, 0666); err != nil {
		fmt.Println("failed to write file:", err)
		return
	}
}
