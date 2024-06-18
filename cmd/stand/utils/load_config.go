package utils

import (
	"github.com/ischenkx/kantoku/pkg/lib/platform"
	"github.com/joho/godotenv"
	"log"
)

func LoadConfig() platform.Config {
	if err := godotenv.Load("cmd/stand/host.env"); err != nil {
		panic(err)
	}

	cfg, err := platform.FromFile("cmd/stand/config.yaml")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	return cfg
}
