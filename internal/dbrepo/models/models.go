package models

// Models stores repositories of each data model, provided that it is also added in the New function.
type Models struct {
	User         UserRepository
	Currency     CurrencyRepository
	Subscription SubscriptionRepository
}
