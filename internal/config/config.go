package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	GitHubToken      string
	GitHubUsername   string
	OpengistURL      string
	OpengistUsername string
	OpengistToken    string
	SyncInterval     time.Duration
	WorkDir          string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// Ignore error if .env file doesn't exist
	}

	config := &Config{
		GitHubToken:      getEnv("GITHUB_TOKEN", ""),
		GitHubUsername:   getEnv("GITHUB_USERNAME", ""),
		OpengistURL:      getEnv("OPENGIST_URL", "http://localhost:6157"),
		OpengistUsername: getEnv("OPENGIST_USERNAME", ""),
		OpengistToken:    getEnv("OPENGIST_TOKEN", ""),
		WorkDir:          getEnv("WORK_DIR", "/tmp/gist-sync"),
		SyncInterval:     getDurationEnv("SYNC_INTERVAL_MINUTES", 30),
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.GitHubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN is required")
	}
	if c.GitHubUsername == "" {
		return fmt.Errorf("GITHUB_USERNAME is required")
	}
	if c.OpengistURL == "" {
		return fmt.Errorf("OPENGIST_URL is required")
	}
	if c.OpengistUsername == "" {
		return fmt.Errorf("OPENGIST_USERNAME is required")
	}
	if c.OpengistToken == "" {
		return fmt.Errorf("OPENGIST_TOKEN is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultMinutes int) time.Duration {
	if value := os.Getenv(key); value != "" {
		if minutes, err := strconv.Atoi(value); err == nil {
			return time.Duration(minutes) * time.Minute
		}
	}
	return time.Duration(defaultMinutes) * time.Minute
}
