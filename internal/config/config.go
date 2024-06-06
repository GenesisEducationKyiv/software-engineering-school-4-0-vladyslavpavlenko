package config

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"gorm.io/gorm"
)

// AppConfig holds the application config.
type AppConfig struct {
	DB          *gorm.DB
	Models      models.Models
	EmailConfig email.Config
}
