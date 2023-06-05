package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug              bool
	DatebaseConnection string
	Domain             string
	Origin             string
	JwtSecret          string
}

func New(path string) (Config, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.NewDecoder(configFile).Decode(&config); err != nil {
		return Config{}, fmt.Errorf("failed to decode config file: %w", err)
	}
	return config, nil
}
