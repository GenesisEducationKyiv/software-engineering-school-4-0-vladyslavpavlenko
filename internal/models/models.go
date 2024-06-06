package models

import (
	"gorm.io/gorm"
)

var db *gorm.DB

// Models is the type for this package. Any model that is included as a member in this type is available throughout
// the application, anywhere that the app variable is used, provided that the model is also added in the New function.
type Models struct {
	User         User
	Currency     Currency
	Subscription Subscription
}

// New is the function used to create an instance of the models package. It returns the type Model, which embeds
// all the types we want to be available to our application.
func New(dbPool *gorm.DB) Models {
	db = dbPool

	return Models{
		User:         User{},
		Currency:     Currency{},
		Subscription: Subscription{},
	}
}
