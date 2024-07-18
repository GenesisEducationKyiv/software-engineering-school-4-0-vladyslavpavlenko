package handlers

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/notifier"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi"
)

type (
	// Services is the repository type.
	Services struct {
		Fetcher    rateapi.Fetcher
		Notifier   *notifier.Notifier
		Subscriber subscriber
	}

	// Repository is the repository type
	Repository struct {
		App      *config.AppConfig
		Services *Services
	}

	// subscriber defines an interface for managing subscriptions.
	subscriber interface {
		AddSubscription(emailAddr string) error
		DeleteSubscription(emailAddr string) error
		GetSubscriptions(limit, offset int) ([]models.Subscription, error)
	}
)

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig, services *Services) *Repository {
	return &Repository{
		App:      a,
		Services: services,
	}
}

// NewHandlers sets the Repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}
