package config

type Config struct {
	System   SystemConfig   `json:"system" yaml:"system"`
	Services ServicesConfig `json:"services" yaml:"services"`
}

type SystemConfig struct {
	Tasks     DynamicConfig `json:"tasks" yaml:"tasks"`
	Resources DynamicConfig `json:"resources" yaml:"resources"`
	Events    DynamicConfig `json:"events" yaml:"events"`
}

type ServicesConfig struct {
	Scheduler  DynamicConfig `yaml:"scheduler"`
	Processor  DynamicConfig `yaml:"processor"`
	Discovery  DynamicConfig `yaml:"discovery"`
	HttpServer DynamicConfig `yaml:"http_api"`
	Status     DynamicConfig `yaml:"status"`
}
