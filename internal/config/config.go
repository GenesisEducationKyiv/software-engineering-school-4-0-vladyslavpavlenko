package config

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"gorm.io/gorm"
)

// AppConfig holds the application config.
type AppConfig struct {
	DB          *gorm.DB
	Models      dbrepo.Models
	EmailConfig email.Config
}
