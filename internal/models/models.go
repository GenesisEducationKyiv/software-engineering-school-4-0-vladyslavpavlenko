package models

import (
	"gorm.io/gorm"
)

// Models stores repositories of each data model, provided that it is also added in the New function.
type Models struct {
	User         *UserRepository
	Currency     *CurrencyRepository
	Subscription *SubscriptionRepository
}

// New creates and returns an instance of Models, which embeds all the models' repositories.
func New(db *gorm.DB) Models {
	return Models{
		User:         NewUserRepository(db),
		Currency:     NewCurrencyRepository(db),
		Subscription: NewSubscriptionRepository(db),
	}
}
