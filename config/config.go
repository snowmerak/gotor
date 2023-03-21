package config

type Config struct {
	Actors      []Actor        `yaml:"actors"`
	Directories map[string]any `yaml:"directories"`
}

type Channel struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type Actor struct {
	Path        string    `yaml:"path"`
	PackageName string    `yaml:"package_name"`
	ActorName   string    `yaml:"actor_name"`
	Channels    []Channel `yaml:"channels"`
}
