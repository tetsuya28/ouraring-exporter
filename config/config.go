package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	Port                           string `env:"PORT" envDefault:"8080"`
	ShutdownTimeout                int64  `env:"SHUTDOWN_TIMEOUT" envDefault:"5"`
	OuraringAPIKey                 string `env:"OURARING_API_KEY,required"`
	OuraringAPICallIntervalSeconds int64  `env:"OURARING_API_CALL_INTERVAL_SECONDS" envDefault:"60"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
