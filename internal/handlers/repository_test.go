package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
)

// MockFetcher is a mock type for the Fetcher interface
type MockFetcher struct {
	mock.Mock
}

func (m *MockFetcher) Fetch() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// MockSubscription is a mock type for the Subscription interface
type MockSubscriber struct {
	mock.Mock
}

func (m *MockSubscriber) AddSubscription(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSubscriber) GetSubscriptions() ([]models.Subscription, error) {
	args := m.Called()
	return args.Get(0).([]models.Subscription), args.Error(1)
}

func TestNewRepo(t *testing.T) {
	appConfig := &config.AppConfig{}
	mockFetcher := new(MockFetcher)
	mockSubscriber := new(MockSubscriber)
	mockSender := new(MockEmailSender)

	repo := handlers.NewRepo(appConfig, mockFetcher, mockSubscriber, mockSender)
	if repo.App != appConfig || repo.Fetcher != mockFetcher || repo.Subscriber != mockSubscriber {
		t.Errorf("NewRepo did not initialize the repository correctly")
	}
}

func TestNewHandlers(t *testing.T) {
	appConfig := &config.AppConfig{}
	mockFetcher := new(MockFetcher)
	mockSubscriber := new(MockSubscriber)
	mockSender := new(MockEmailSender)

	repo := handlers.NewRepo(appConfig, mockFetcher, mockSubscriber, mockSender)
	handlers.NewHandlers(repo)
	if handlers.Repo != repo {
		t.Errorf("NewHandlers did not set the global Repo variable correctly")
	}
}
