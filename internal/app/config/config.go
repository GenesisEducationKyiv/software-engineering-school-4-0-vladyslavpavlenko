package config

import (
	"github.com/vladyslavpavlenko/genesis-api-project/pkg/logger"
)

// Config holds the application config.
type Config struct {
	Logger *logger.Logger
}

// New creates a new Config.
func New() *Config {
	return &Config{}
}
