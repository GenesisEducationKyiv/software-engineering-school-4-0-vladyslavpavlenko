package config

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox"
)

// AppConfig holds the application config.
type AppConfig struct {
	Outbox outbox.Outbox
}

// NewAppConfig creates a new AppConfig.
func NewAppConfig() *AppConfig {
	return &AppConfig{}
}
