package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
)

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) Send(cfg email.Config, params email.Params) error {
	args := m.Called(cfg, params)
	return args.Error(0)
}

func TestNotifySubscribers_Success(t *testing.T) {
	mockSubscriber := new(MockSubscriber)
	mockFetcher := new(MockFetcher)
	mockSender := new(MockEmailSender)
	appConfig := &config.AppConfig{}
	repo := handlers.NewRepo(appConfig, mockFetcher, mockSubscriber, mockSender)

	subscribers := []models.Subscription{{Email: "user@example.com"}}
	mockSubscriber.On("GetSubscriptions").Return(subscribers, nil)
	mockFetcher.On("Fetch").Return("24.5", nil)
	mockSender.On("Send", mock.AnythingOfType("email.Config"), mock.AnythingOfType("email.Params")).Return(nil)

	repo.Sender = mockSender
	err := repo.NotifySubscribers()
	assert.NoError(t, err)

	mockSubscriber.AssertExpectations(t)
	mockFetcher.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}
