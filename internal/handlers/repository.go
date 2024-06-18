package handlers

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App        *config.AppConfig
	Subscriber Subscriber
	Fetcher    Fetcher
	Sender     Sender
}

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig, fetcher Fetcher, subscriber Subscriber, sender Sender) *Repository {
	return &Repository{
		App:        a,
		Subscriber: subscriber,
		Fetcher:    fetcher,
		Sender:     sender,
	}
}

// NewHandlers sets the Repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}
