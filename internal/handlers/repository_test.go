package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
)

// MockDB is a mock implementation of the dbrepo.DB interface
type MockDB struct{}

func (m *MockDB) Connect(_ string) error {
	return nil
}

func (m *MockDB) Close() error {
	return nil
}

func (m *MockDB) Migrate() error {
	return nil
}

func (m *MockDB) SomeDBFunction() error {
	return nil
}

func TestNewRepo(t *testing.T) {
	appConfig := &config.AppConfig{}
	mockDB := &MockDB{}

	repo := handlers.NewRepo(appConfig, mockDB)
	assert.Equal(t, appConfig, repo.App, "AppConfig should be correctly assigned in NewRepo")
	assert.Equal(t, mockDB, repo.DB, "DB should be correctly assigned in NewRepo")
}

func TestNewHandlers(t *testing.T) {
	appConfig := &config.AppConfig{}
	mockDB := &MockDB{}

	repo := handlers.NewRepo(appConfig, mockDB)
	handlers.NewHandlers(repo)
	assert.Equal(t, repo, handlers.Repo, "Repo should be correctly set by NewHandlers")
}
