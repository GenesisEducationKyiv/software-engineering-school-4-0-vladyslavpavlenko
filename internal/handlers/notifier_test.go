package handlers_test

import (
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/storage/gormrepo"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
)

func TestNotifySubscribers_Success(t *testing.T) {
	mockDB := new(MockDB)
	mockFetcher := new(MockFetcher)
	appConfig := &config.AppConfig{}
	dbConn := &gormrepo.Connection{}
	services := setupServicesWithMocks(mockFetcher)
	repo := handlers.NewRepo(appConfig, services, dbConn)

	subscribers := []models.Subscription{{Email: "user@example.com"}}
	mockDB.On("GetSubscriptions").Return(subscribers, nil)
	mockFetcher.On("Fetch", mock.Anything, "USD", "UAH").Return("24.5", nil)

	err := repo.ProduceMailingEvents()

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
	mockFetcher.AssertExpectations(t)
}
