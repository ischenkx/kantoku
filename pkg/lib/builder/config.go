package builder

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
	"time"
)

type SystemConfig struct {
	Tasks     TasksConfig     `yaml:"tasks,omitempty" json:"tasks,omitempty"`
	Resources ResourcesConfig `yaml:"resources,omitempty" json:"resources,omitempty"`
	Events    EventsConfig    `yaml:"events,omitempty" json:"events,omitempty"`
}

type TasksConfig struct {
	Storage TasksStorageConfig `yaml:"storage,omitempty" json:"storage,omitempty"`
}

type TasksStorageConfig struct {
	Kind    string         `yaml:"kind,omitempty" json:"kind,omitempty"`
	URI     string         `yaml:"uri,omitempty" json:"uri,omitempty"`
	Options map[string]any `yaml:"options,omitempty" json:"options,omitempty"`
}

type ResourcesConfig struct {
	Storage   ResourcesStorageConfig    `yaml:"storage,omitempty" json:"storage,omitempty"`
	Observers []ResourcesObserverConfig `yaml:"observers,omitempty" json:"observers,omitempty"`
}

type ResourcesStorageConfig struct {
	Kind    string         `yaml:"kind,omitempty" json:"kind,omitempty"`
	URI     string         `yaml:"uri,omitempty" json:"uri,omitempty"`
	Options map[string]any `yaml:"options,omitempty" json:"options,omitempty"`
}

type ResourcesObserverConfig struct {
	Kind    string        `yaml:"kind,omitempty" json:"kind,omitempty"`
	Options DynamicConfig `yaml:"options,omitempty" json:"options,omitempty"`
}

type EventsConfig struct {
	Broker EventsBrokerConfig `yaml:"broker,omitempty" json:"broker,omitempty"`
}

type EventsBrokerConfig struct {
	Kind    string         `yaml:"kind,omitempty" json:"kind,omitempty"`
	URI     string         `yaml:"uri,omitempty" json:"uri,omitempty"`
	Options map[string]any `yaml:"options,omitempty" json:"options,omitempty"`
}

type SpecificationsConfig struct {
	Storage SpecificationsStorageConfig `yaml:"storage,omitempty" json:"storage,omitempty"`
}

type SpecificationsStorageConfig struct {
	Kind    string         `yaml:"kind,omitempty" json:"kind,omitempty"`
	URI     string         `yaml:"uri,omitempty" json:"uri,omitempty"`
	Options map[string]any `yaml:"options,omitempty" json:"options,omitempty"`
}

type HttpApiServiceConfig struct {
	ServiceConfig ServiceConfig `yaml:"$,omitempty" json:"$,omitempty"`
	Port          int           `yaml:"port,omitempty" json:"port,omitempty"`
	LoggerEnabled bool          `yaml:"logger_enabled,omitempty" json:"logger_enabled,omitempty"`
}

type SchedulerServiceConfig struct {
	ServiceConfig ServiceConfig               `yaml:"$,omitempty" json:"$,omitempty"`
	TaskToGroup   SchedulerTaskToGroupConfig  `yaml:"task_to_group,omitempty" json:"task_to_group,omitempty"`
	Dependencies  SchedulerDependenciesConfig `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
	Resolvers     []SchedulerResolverConfig   `yaml:"resolvers,omitempty" json:"resolvers,omitempty"`
}

type SchedulerTaskToGroupConfig struct {
	Kind    string         `yaml:"kind,omitempty" json:"kind,omitempty"`
	URI     string         `yaml:"uri,omitempty" json:"uri,omitempty"`
	Options map[string]any `yaml:"options,omitempty" json:"options,omitempty"`
}

type SchedulerDependenciesConfig struct {
	Poller struct {
		Interval  time.Duration `yaml:"interval,omitempty" json:"interval,omitempty"`
		BatchSize int           `yaml:"batch_size,omitempty" json:"batch_size,omitempty"`
	} `yaml:"poller,omitempty" json:"poller,omitempty"`
	Postgres struct {
		URI     string         `yaml:"uri,omitempty" json:"uri,omitempty"`
		Options map[string]any `yaml:"options,omitempty" json:"options,omitempty"`
	} `yaml:"postgres,omitempty" json:"postgres,omitempty"`
}

type SchedulerResolverConfig struct {
	Kind string        `yaml:"kind,omitempty" json:"kind,omitempty"`
	Data DynamicConfig `yaml:"data,omitempty" json:"data,omitempty"`
}

type SchedulerResourceResolverConfig struct {
	Storage SchedulerResourceResolverStorageConfig `yaml:"storage,omitempty" json:"storage,omitempty"`
	Poller  struct {
		Limit    int           `yaml:"limit,omitempty" json:"limit,omitempty"`
		Interval time.Duration `yaml:"interval,omitempty" json:"interval,omitempty"`
	} `yaml:"poller,omitempty" json:"poller,omitempty"`
	//Topic string `yaml:"topic,omitempty" json:"topic,omitempty"`
}

type SchedulerResourceResolverStorageConfig struct {
	Kind    string         `yaml:"kind,omitempty" json:"kind,omitempty"`
	URI     string         `yaml:"uri,omitempty" json:"uri,omitempty"`
	Options map[string]any `yaml:"options,omitempty" json:"options,omitempty"`
}

type StatusServiceConfig struct {
	ServiceConfig ServiceConfig `yaml:"$,omitempty" json:"$,omitempty"`
}

type ProcessorServiceConfig struct {
	ServiceConfig ServiceConfig `yaml:"$,omitempty" json:"$,omitempty"`
	Kind          string        `yaml:"kind,omitempty" json:"kind,omitempty"`
}

type DiscoveryServiceConfig struct {
	ServiceConfig   ServiceConfig      `yaml:"$,omitempty" json:"$,omitempty"`
	PollingInterval time.Duration      `yaml:"polling_interval,omitempty" json:"polling_interval,omitempty"`
	Hub             DiscoveryHubConfig `yaml:"hub,omitempty" json:"hub,omitempty"`
}

type DiscoveryHubConfig struct {
	Kind    string         `yaml:"kind,omitempty" json:"kind,omitempty"`
	Options map[string]any `yaml:"options,omitempty" json:"options,omitempty"`
}

type ServiceConfig struct {
	ID        string                 `yaml:"id,omitempty" json:"id,omitempty"`
	Name      string                 `yaml:"name,omitempty" json:"name,omitempty"`
	Discovery ServiceDiscoveryConfig `yaml:"discovery,omitempty" json:"discovery,omitempty"`
}

type ServiceDiscoveryConfig struct {
	Enabled bool           `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Info    map[string]any `yaml:"info,omitempty" json:"info,omitempty"`
}

type Config struct {
	Core     CoreConfig     `yaml:"core,omitempty" json:"core,omitempty"`
	Services ServicesConfig `yaml:"services,omitempty" json:"services,omitempty"`
}

type CoreConfig struct {
	System         SystemConfig         `yaml:"system,omitempty" json:"system,omitempty"`
	Specifications SpecificationsConfig `yaml:"specifications,omitempty" json:"specifications,omitempty"`
}

type ServicesConfig struct {
	HttpApi   HttpApiServiceConfig   `yaml:"http_api,omitempty" json:"http_api,omitempty"`
	Scheduler SchedulerServiceConfig `yaml:"scheduler,omitempty" json:"scheduler,omitempty"`
	Status    StatusServiceConfig    `yaml:"status,omitempty" json:"status,omitempty"`
	Processor ProcessorServiceConfig `yaml:"processor,omitempty" json:"processor,omitempty"`
	Discovery DiscoveryServiceConfig `yaml:"discovery,omitempty" json:"discovery,omitempty"`
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
			data, err := io.ReadAll(from)
			if err != nil {
				return fmt.Errorf("failed to read from file: %w", err)
			}

			stringData := string(data)

			expandedData := os.ExpandEnv(stringData)

			return yaml.NewDecoder(strings.NewReader(expandedData)).Decode(to)
		}, nil
	case strings.HasSuffix(path, ".json"):
		return func(from io.Reader, to any) error {
			data, err := io.ReadAll(from)
			if err != nil {
				return fmt.Errorf("failed to read from file: %w", err)
			}

			stringData := string(data)

			expandedData := os.ExpandEnv(stringData)

			return json.NewDecoder(strings.NewReader(expandedData)).Decode(to)
		}, nil
	default:
		return nil, fmt.Errorf("unknown format")
	}
}
