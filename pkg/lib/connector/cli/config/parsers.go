package config

import (
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
)

func FromEnv(prefix string) (config Config, err error) {
	err = envconfig.Process(prefix, &config)
	return
}

func FromFile(path string) (config Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open the file: %w", err)
	}
	defer file.Close()

	decoder, err := decoderByPath(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to make a decoder: %w", err)
	}

	if err := decoder(file, &config); err != nil {
		return Config{}, fmt.Errorf("failed to decode: %w", err)
	}

	return
}

func decoderByPath(path string) (decoder func(from io.Reader, to any) error, err error) {
	path = strings.TrimSpace(path)

	switch {
	case strings.HasSuffix(path, ".yaml"), strings.HasSuffix(path, ".yml"):
		return func(from io.Reader, to any) error {
			return yaml.NewDecoder(from).Decode(to)
		}, nil
	case strings.HasSuffix(path, ".json"):
		return func(from io.Reader, to any) error {
			return json.NewDecoder(from).Decode(to)
		}, nil
	default:
		return nil, fmt.Errorf("unknown format")
	}
}
