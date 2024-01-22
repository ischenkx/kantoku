package cli

type Config struct {
	System SystemConfig `json:"system" yaml:"system"`

	Services ServicesConfig
}

type SystemConfig struct {
	Tasks     TasksConfig     `json:"tasks" yaml:"tasks"`
	Resources ResourcesConfig `json:"resources" yaml:"resources"`
	Events    EventsConfig    `json:"events" yaml:"events"`
}

type TasksConfig struct {
	Addr    string         `json:"addr" yaml:"addr"`
	Options map[string]any `json:"options" yaml:"options"`
}

type ResourcesConfig struct {
	Addr string `json:"addr" yaml:"addr"`
}

type EventsConfig struct {
	Addr string `json:"addr" yaml:"addr"`
}

type ServicesConfig struct {
	Scheduler SchedulerConfig
}

type SchedulerConfig struct {
	Type    string
	Options map[string]any
}
