package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Jaeger  string `envconfig:"JAEGER" default:"http://localhost:14268/api/traces"`
	Tracing bool   `envconfig:"TRACING" default:"false"`
}

var config *Config
var once sync.Once

// Load reads config file and ENV variables if set.
func Load() *Config {
	once.Do(func() {
		load()
	})

	return config
}

func load() {
	config = new(Config)
	if err := envconfig.Process("", config); err != nil {
		log.Fatal().Err(err).Send()
	}
}
