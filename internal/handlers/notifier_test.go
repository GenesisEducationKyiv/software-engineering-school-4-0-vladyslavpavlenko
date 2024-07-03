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
	mockSubscriber := new(MockSubscriber)
	mockFetcher := new(MockFetcher)
	mockSender := new(MockSender)
	appConfig := &config.AppConfig{}
	dbConn := &gormrepo.Connection{}
	services := setupServicesWithMocks(mockFetcher, mockSender)
	repo := handlers.NewRepo(appConfig, services, dbConn)

	subscribers := []models.Subscription{{Email: "user@example.com"}}
	mockSubscriber.On("GetSubscriptions").Return(subscribers, nil)
	mockFetcher.On("Fetch", mock.Anything, "USD", "UAH").Return("24.5", nil)
	mockSender.On("Send", mock.AnythingOfType("email.Config"), mock.AnythingOfType("email.Params")).Return(nil)

	err := repo.NotifySubscribers()

	assert.NoError(t, err)

	mockSubscriber.AssertExpectations(t)
	mockFetcher.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}
