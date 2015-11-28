package gamq

type Config struct {
	Port             int
	ProfilingEnabled bool
	StatsDEndpoint   string
}

var Configuration *Config

func SetConfig(givenConfig *Config) {
	Configuration = givenConfig
}
