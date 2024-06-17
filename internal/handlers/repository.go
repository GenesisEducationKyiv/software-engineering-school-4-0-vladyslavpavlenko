package handlers

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App          *config.AppConfig
	Subscription Subscription
	Fetcher      Fetcher
}

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig, fetcher Fetcher, subscription Subscription) *Repository {
	return &Repository{
		App:          a,
		Subscription: subscription,
		Fetcher:      fetcher,
	}
}

// NewHandlers sets the Repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}
