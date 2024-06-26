package config

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
)

// AppConfig holds the application config.
type AppConfig struct {
	EmailConfig email.Config
}

// NewAppConfig creates a new AppConfig.
func NewAppConfig() *AppConfig {
	return &AppConfig{EmailConfig: email.Config{}}
}
