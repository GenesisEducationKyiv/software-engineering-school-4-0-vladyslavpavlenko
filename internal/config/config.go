package config

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
)

// AppConfig holds the application config.
type AppConfig struct {
	Models      models.Models
	EmailConfig email.Config
}
