package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/impl/broker/watermill"
	redisResources "github.com/ischenkx/kantoku/pkg/lib/impl/core/resource/redis"
	mongorec "github.com/ischenkx/kantoku/pkg/lib/impl/data/record/mongo"
	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
	"io"
	"net/url"
	"os"
	"strings"
)

func configFromEnv() (config Config, err error) {
	err = envconfig.Process("KANTOKU", &config)
	return
}

func configFromFile(path string) (config Config, err error) {
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

func systemFromConfig(config SystemConfig) (system.AbstractSystem, error) {
	events, err := eventsFromConfig(config.Events)
	if err != nil {
		return nil, fmt.Errorf("failed to create events: %w", err)
	}

	resources, err := resourcesFromConfig(config.Resources)
	if err != nil {
		return nil, fmt.Errorf("failed to create resources: %w", err)
	}

	tasks, err := tasksFromConfig(config.Tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to create tasks: %w", err)
	}

	return system.New(events, resources, tasks), nil
}

func eventsFromConfig(config EventsConfig) (*event.Broker, error) {
	u, err := url.Parse(config.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	var agent watermill.Agent

	switch u.Scheme {
	case "nats":
		agent, err = watermill.Nats(config.Addr, nats.SubscriberConfig{}, nats.PublisherConfig{})
		if err != nil {
			return nil, fmt.Errorf("failed to create a nats agent: %w", err)
		}
	case "redis":
		client := redis.NewClient(&redis.Options{
			Addr: config.Addr,
		})
		agent, err = watermill.Redis(client, redisstream.SubscriberConfig{}, redisstream.PublisherConfig{})
		if err != nil {
			return nil, fmt.Errorf("failed to create a nats agent: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported agent: %s", u.Scheme)
	}

	baseBroker := watermill.Broker[event.Event]{
		Agent:                     agent,
		ItemCodec:                 codec.JSON[event.Event](),
		ConsumerChannelBufferSize: 1024,
	}

	return event.NewBroker(baseBroker), nil
}

func resourcesFromConfig(config ResourcesConfig) (resource.Storage, error) {
	u, err := url.Parse(config.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	switch u.Scheme {
	case "redis":
		client := redis.NewClient(&redis.Options{
			Addr: config.Addr,
		})
		storage := redisResources.New(client, codec.JSON[resource.Resource](), "resources")

		return storage, nil
	default:
		return nil, fmt.Errorf("unsupported agent: %s", u.Scheme)
	}
}

func tasksFromConfig(config TasksConfig) (record.Storage[task.Task], error) {
	u, err := url.Parse(config.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	switch u.Scheme {
	case "mongo":
		clientOptions := options.Client().ApplyURI(config.Addr)

		// Connect to the MongoDB server
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to mongo: %w", err)
		}

		// TODO use config.Options + mapstructure to use dynamic configs

		storage := mongorec.New[task.Task](
			client.Database("main").Collection("tasks"),
			task.Codec{},
		)

		return storage, nil
	default:
		return nil, fmt.Errorf("unsupported agent: %s", u.Scheme)
	}
}
