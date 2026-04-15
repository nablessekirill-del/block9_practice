package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AppName     string `yaml:"app_name"`
	Currency    string `yaml:"currency"`
	MaxProducts int    `yaml:"max_products"`
	LogLevel    string `yaml:"log_level"`
}

func Load(path string) (*Config, error) {
	// дефолтные значения — если поля нет в yml, будет это
	cfg := &Config{
		LogLevel: "info",
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
