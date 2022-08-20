package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DB                   string `env:"DATABASE_URI"`
}

func NewConfig() (config, error) {
	cfg := config{}
	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "set server address, by example: 127.0.0.1:8080")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "set accrual server address, by example: http://127.0.0.1")
	flag.StringVar(&cfg.DB, "d", cfg.DB, "set database URI for Postgres, by example: host=localhost port=5432 user=example password=123 dbname=example sslmode=disable connect_timeout=5")

	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	flag.Parse()
	return cfg, nil
}
