package gormrepo

import "github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/models"

// NewModels creates and returns an instance of models.Models, which embeds all the models' repositories.
func NewModels(conn *GormDB) models.Models {
	return models.Models{
		User:         NewGormUserRepository(conn),
		Currency:     NewGormCurrencyRepository(conn),
		Subscription: NewGormSubscriptionRepository(conn),
	}
}
