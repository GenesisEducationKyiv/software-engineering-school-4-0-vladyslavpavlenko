package handlers_test

import (
	"context"
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"

	"github.com/stretchr/testify/assert"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/notifier"
)

type mockFetcher struct{}

func (m *mockFetcher) Fetch(_ context.Context, _, _ string) (string, error) {
	return "mocked-fetch", nil
}

type mockSubscriber struct{}

func (m *mockSubscriber) AddSubscription(_ string) error {
	return nil
}

func (m *mockSubscriber) DeleteSubscription(_ string) error {
	return nil
}

func (m *mockSubscriber) GetSubscriptions(_, _ int) ([]models.Subscription, error) {
	return []models.Subscription{}, nil
}

func TestNewHandlers(t *testing.T) {
	appConfig := &config.Config{}

	mockServices := &handlers.Services{
		Fetcher:    &mockFetcher{},
		Notifier:   &notifier.Notifier{},
		Subscriber: &mockSubscriber{},
	}

	h := handlers.NewHandlers(appConfig, mockServices)

	assert.NotNil(t, h)
	assert.Equal(t, appConfig, h.App)
	assert.Equal(t, mockServices, h.Services)
}
