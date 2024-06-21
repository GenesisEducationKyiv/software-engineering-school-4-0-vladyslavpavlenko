package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
)

// TestNewRepo tests the creation of a new repository
func TestNewRepo(t *testing.T) {
	appConfig := &config.AppConfig{}
	services := &handlers.Services{
		Subscriber: &MockSubscriber{},
		Fetcher:    &MockFetcher{},
		Sender:     &MockSender{},
	}

	repo := handlers.NewRepo(appConfig, services)

	assert.NotNil(t, repo)
	assert.Equal(t, appConfig, repo.App)
	assert.Equal(t, services, repo.Services)
}

// TestNewHandlers tests setting the repository
func TestNewHandlers(t *testing.T) {
	appConfig := &config.AppConfig{}
	services := &handlers.Services{
		Subscriber: &MockSubscriber{},
		Fetcher:    &MockFetcher{},
		Sender:     &MockSender{},
	}

	repo := handlers.NewRepo(appConfig, services)
	handlers.NewHandlers(repo)

	assert.Equal(t, repo, handlers.Repo)
}
