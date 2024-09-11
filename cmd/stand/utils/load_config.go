package utils

import (
	"github.com/ischenkx/kantoku/pkg/lib/builder"
	"github.com/joho/godotenv"
	"log"
)

func LoadConfig() builder.Config {
	if err := godotenv.Load("cmd/stand/host.env"); err != nil {
		panic(err)
	}

	cfg, err := builder.FromFile("cmd/stand/config.yaml")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	return cfg
}
