package gormrepo

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo"
	"gorm.io/gorm"
)

// New creates and returns an instance of Models, which embeds all the models' repositories.
func New(db *gorm.DB) dbrepo.Models {
	return dbrepo.Models{
		User:         NewGormUserRepository(db),
		Currency:     NewGormCurrencyRepository(db),
		Subscription: NewGormSubscriptionRepository(db),
	}
}
