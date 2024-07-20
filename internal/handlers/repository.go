package handlers

import (
	"context"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/notifier"
)

type fetcher interface {
	Fetch(ctx context.Context, base, target string) (string, error)
}

// Services is the repository type.
type Services struct {
	Fetcher    fetcher
	Notifier   *notifier.Notifier
	Subscriber subscriber
}

// Repository is the repository type
type Repository struct {
	App      *config.Config
	Services *Services
}

// subscriber defines an interface for managing subscriptions.
type subscriber interface {
	AddSubscription(emailAddr string) error
	DeleteSubscription(emailAddr string) error
	GetSubscriptions(limit, offset int) ([]models.Subscription, error)
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new Repository
func NewRepo(a *config.Config, services *Services) *Repository {
	return &Repository{
		App:      a,
		Services: services,
	}
}

// NewHandlers sets the Repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}
