package handlers

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi"
)

type (
	// Services is the repository type.
	Services struct {
		Fetcher rateapi.Fetcher
	}

	// Repository is the repository type
	Repository struct {
		App      *config.AppConfig
		DB       dbConnection
		Services *Services
	}

	// dbConnection defines an interface for the database connection.
	dbConnection interface {
		AddSubscription(emailAddr string) error
		GetSubscriptions(limit, offset int) ([]models.Subscription, error)
	}
)

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig, services *Services, conn dbConnection) *Repository {
	return &Repository{
		App:      a,
		DB:       conn,
		Services: services,
	}
}

// NewHandlers sets the Repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}
