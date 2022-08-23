package config

type Config struct {
}

var (
	config *Config
)

func Init(c *Config) {
	config = c
}

func Get() *Config {
	return config
}
