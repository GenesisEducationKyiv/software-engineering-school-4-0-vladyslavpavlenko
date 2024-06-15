package config

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
)

// AppConfig holds the application config.
type AppConfig struct {
	DB          dbrepo.DB
	Models      models.Models
	EmailConfig email.Config
}
