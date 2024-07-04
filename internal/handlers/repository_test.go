package handlers_test

import (
	"context"
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/storage/gormrepo"

	"github.com/stretchr/testify/mock"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"

	"github.com/stretchr/testify/assert"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
)

type MockSender struct {
	mock.Mock
}

func (m *MockSender) Send(cfg email.Config, params email.Params) error {
	args := m.Called(cfg, params)
	return args.Error(0)
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetSubscriptions(_, _ int) ([]models.Subscription, error) {
	args := m.Called()
	return args.Get(0).([]models.Subscription), args.Error(1)
}

func (m *MockDB) AddSubscription(emailAddr string) error {
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

func setupServicesWithMocks(fetcher *MockFetcher) *handlers.Services {
	return &handlers.Services{
		Fetcher: fetcher,
	}
}

// TestNewRepo tests the creation of a new repository
func TestNewRepo(t *testing.T) {
	appConfig := &config.AppConfig{}
	services := &handlers.Services{
		Fetcher: &MockFetcher{},
	}
	dbConn := &gormrepo.Connection{}

	repo := handlers.NewRepo(appConfig, services, dbConn)

	assert.NotNil(t, repo)
	assert.Equal(t, appConfig, repo.App)
	assert.Equal(t, services, repo.Services)
}

// TestNewHandlers tests setting the repository
func TestNewHandlers(t *testing.T) {
	appConfig := &config.AppConfig{}
	services := &handlers.Services{
		Fetcher: &MockFetcher{},
	}
	dbConn := &gormrepo.Connection{}

	repo := handlers.NewRepo(appConfig, services, dbConn)
	handlers.NewHandlers(repo)

	assert.Equal(t, repo, handlers.Repo)
}
