package config

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/notifier"
)

// AppConfig holds the application config.
type AppConfig struct {
	Notifier notifier.Notifier
}

// NewAppConfig creates a new AppConfig.
func NewAppConfig() *AppConfig {
	return &AppConfig{}
}
