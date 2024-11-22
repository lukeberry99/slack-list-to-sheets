package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SlackToken string
	ServerPort string
}

func Load() (*Config, error) {
	// Load .env file if it exists, ignore error if it doesn't
	_ = godotenv.Load()

	config := &Config{
		SlackToken: os.Getenv("SLACK_TOKEN"),
		ServerPort: os.Getenv("SERVER_PORT"),
	}

	if config.ServerPort == "" {
		config.ServerPort = "8080"
	}

	if config.SlackToken == "" {
		return nil, fmt.Errorf("SLACK_TOKEN is not set")
	}

	return config, nil
}
