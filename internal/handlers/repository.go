package handlers

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"
)

type (
	// Services is the repository type.
	Services struct {
		Subscriber Subscriber
		Fetcher    Fetcher
		Sender     Sender
	}

	// Repository is the repository type
	Repository struct {
		App      *config.AppConfig
		Services *Services
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
