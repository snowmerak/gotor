package config

type Config struct {
	Actors []struct {
		Path        string `yaml:"path"`
		PackageName string `yaml:"package_name"`
		ActorName   string `yaml:"actor_name"`
		Channels    []struct {
			Name string `yaml:"name"`
			Type string `yaml:"type"`
		} `yaml:"channels"`
	} `yaml:"actors"`
}
