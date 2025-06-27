package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
)

const (
	defaultRunAddress           = "localhost:8080"
	defaultDatabaseURI          = "postgres://postgres:postgres@localhost:5432/postgres"
	defaultAccrualSystemAddress = "localhost:8081"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func New() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.RunAddress, "a", defaultRunAddress, "run address")
	flag.StringVar(&cfg.DatabaseURI, "d", defaultDatabaseURI, "database uri")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", defaultAccrualSystemAddress, "accrual system address")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
