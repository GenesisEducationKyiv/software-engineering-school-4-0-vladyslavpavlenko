package models

// Models stores repositories of each data model.
type Models struct {
	User         UserRepository
	Currency     CurrencyRepository
	Subscription SubscriptionRepository
}
