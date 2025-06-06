// Package config парсит конфиг из config.yml
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port int `yaml:"port"`
}

// Парсит параметры из config.yml и возвращает структуру Config, ошибку
func Load() (*Config, error) {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	if cfg.Port == 0 {
		cfg.Port = 8080
	}

	return &cfg, nil
}
