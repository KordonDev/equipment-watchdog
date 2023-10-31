package config

import (
	"fmt"
	"os"

	"github.com/cloudflare/cfssl/log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug              bool
	DatabaseConnection string `yaml:"databaseConnection"`
	Domain             string
	Origin             string
	JwtSecret          string `yaml:"jwtSecret"`
}

func New(path string) (Config, error) {
  dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(dir)

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
