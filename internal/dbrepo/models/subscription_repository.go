package models

// SubscriptionRepository interface defines methods to access Subscription data.
type SubscriptionRepository interface {
	Create(uint, uint, uint) (*Subscription, error)
	GetSubscriptions() ([]Subscription, error)
}
