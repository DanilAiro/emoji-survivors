package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBUrl  string
	Secret string
	Port   string
}

func Load() (Config, error) {
	cfg := Config{
		DBUrl:  os.Getenv("DB_URL"),
		Secret: os.Getenv("JWT_SECRET"),
		Port:   os.Getenv("PORT"),
	}

	if cfg.DBUrl == "" {
		return cfg, fmt.Errorf("переменная окружения DB_URL не задана")
	}
	if cfg.Secret == "" {
		return cfg, fmt.Errorf("переменная окружения JWT_SECRET не задана")
	}
	if cfg.Port == "" {
		cfg.Port = "8082"
	}

	return cfg, nil
}
