package config

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox/producer"
)

// AppConfig holds the application config.
type AppConfig struct {
	Outbox producer.Outbox
}

// NewAppConfig creates a new AppConfig.
func NewAppConfig() *AppConfig {
	return &AppConfig{}
}
