package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
)

type MockSender struct {
	mock.Mock
}

func (m *MockSender) Send(cfg email.Config, params email.Params) error {
	args := m.Called(cfg, params)
	return args.Error(0)
}

type MockSubscriber struct {
	mock.Mock
}

func (m *MockSubscriber) GetSubscriptions(_, _ int) ([]models.Subscription, error) {
	args := m.Called()
	return args.Get(0).([]models.Subscription), args.Error(1)
}

func (m *MockSubscriber) AddSubscription(emailAddr string) error {
	args := m.Called(emailAddr)
	return args.Error(0)
}

type MockFetcher struct {
	mock.Mock
}

func (m *MockFetcher) Fetch(ctx context.Context, base, target string) (string, error) {
	args := m.Called(ctx, base, target)
	return args.String(0), args.Error(1)
}

func setupServicesWithMocks(subscriber *MockSubscriber, fetcher *MockFetcher, sender *MockSender) *handlers.Services {
	return &handlers.Services{
		Subscriber: subscriber,
		Fetcher:    fetcher,
		Sender:     sender,
	}
}

func TestSubscribeUser(t *testing.T) {
	mockFetcher := new(MockFetcher)
	mockSubscriber := new(MockSubscriber)
	mockSender := new(MockSender)
	appConfig := &config.AppConfig{}
	services := setupServicesWithMocks(mockSubscriber, mockFetcher, mockSender)
	repo := handlers.NewRepo(appConfig, services)

	tests := []struct {
		name         string
		email        string
		mockResponse error
		wantStatus   int
		wantErr      string
		expectCreate bool
	}{
		{"Invalid Email", "bademail", nil, http.StatusBadRequest, "invalid email", false},
		{"Valid Email Success", "good@email.com", nil, http.StatusAccepted, "", true},
		{"Duplicate Email", "duplicate@email.com", gorm.ErrDuplicatedKey, http.StatusConflict, "already subscribed", true},
		{"DB Error", "error@email.com", errors.New("db error"), http.StatusInternalServerError, "subscription already exists", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectCreate {
				mockSubscriber.On("AddSubscription", tc.email).Return(tc.mockResponse)
			}

			statusCode, err := repo.SubscribeUser(tc.email)
			assert.Equal(t, tc.wantStatus, statusCode)
			if err != nil {
				assert.Contains(t, err.Error(), tc.wantErr)
			}

			mockSubscriber.AssertExpectations(t)
		})
	}
}
