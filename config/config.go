package config

import (
	"github.com/mattn/go-isatty"
	"os"
)

type Config struct {
}

var (
	config *Config
	tty    bool
)

func Init(c *Config) {
	config = c
}

func Get() *Config {
	return config
}

func init() {
	tty = isatty.IsTerminal(os.Stdout.Fd())
}

func IsTty() bool {
	return tty
}
