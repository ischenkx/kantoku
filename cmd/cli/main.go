package main

import (
	"fmt"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli"
	"os"
)

func main() {
	if err := cli.New().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
