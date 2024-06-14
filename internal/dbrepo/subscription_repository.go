package dbrepo

// SubscriptionRepository interface defines methods to access subscription data.
type SubscriptionRepository interface {
	Create(uint, uint, uint) (*Subscription, error)
	GetSubscriptions() ([]Subscription, error)
}
