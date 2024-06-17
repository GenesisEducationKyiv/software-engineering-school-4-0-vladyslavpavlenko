package handlers

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App     *config.AppConfig
	DB      dbrepo.DB
	Fetcher Fetcher
}

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig, db dbrepo.DB, fetcher Fetcher) *Repository {
	return &Repository{
		App:     a,
		DB:      db,
		Fetcher: fetcher,
	}
}

// NewHandlers sets the Repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}
