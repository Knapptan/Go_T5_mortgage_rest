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

// Парсит параметры из указанного файла
func Load(filename string) (*Config, error) {
	// Читаем файл
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	// Записываем в структуру Config
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	// По умолчанию 8080, если порт не задан или задан некорректно (<=0)
	if cfg.Port <= 0 {
		cfg.Port = 8080
	}

	// Проверяем границы порта
	if cfg.Port > 65535 {
		return nil, fmt.Errorf("invalid port %d: must be between 1 and 65535", cfg.Port)
	}

	return &cfg, nil
}
