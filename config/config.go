package config

type Config struct {
	User  string
	Token string
	Debug bool
}

var (
	config Config
)

func Init(c Config) {
	config = c
}

func Get() Config {
	return config
}
