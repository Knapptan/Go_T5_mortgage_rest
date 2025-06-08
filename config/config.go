// Package config предоставляет загрузку конфигурации приложения из YAML-файла.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config представляет структуру конфигурации приложения.
type Config struct {
	Port int `yaml:"port"`
}

// ErrInvalidPort указывает на недопустимое значение порта.
var ErrInvalidPort = errors.New("port must be between 1 and 65535")

// Load загружает и парсит конфигурацию из указанного YAML-файла.
//
// Если порт не задан или задан некорректно (<= 0), используется порт 8080.
// Если порт превышает 65535, возвращается ошибка.
func Load(filename string) (*Config, error) {
	cleanPath := filepath.Clean(filename)
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	if cfg.Port <= 0 {
		cfg.Port = 8080
	}

	if cfg.Port > 65535 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidPort, cfg.Port)
	}

	return &cfg, nil
}
