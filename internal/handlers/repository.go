package handlers

import (
	"context"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/notifier"
)

type (
	fetcher interface {
		Fetch(ctx context.Context, base, target string) (string, error)
	}

	subscriber interface {
		AddSubscription(emailAddr string) error
		DeleteSubscription(emailAddr string) error
		GetSubscriptions(limit, offset int) ([]models.Subscription, error)
	}
)

// Services is the repository type for the services necessary for API handlers.
type Services struct {
	Fetcher    fetcher
	Notifier   *notifier.Notifier
	Subscriber subscriber
}

// Handlers is the repository type for API handlers.
type Handlers struct {
	App      *config.Config
	Services *Services
}

// NewHandlers creates new Handlers.
func NewHandlers(a *config.Config, services *Services) *Handlers {
	return &Handlers{
		App:      a,
		Services: services,
	}
}
