package handlers_test

import (
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
)

func TestNewRepo(t *testing.T) {
	appConfig := &config.AppConfig{}
	repo := handlers.NewRepo(appConfig)
	if repo.App != appConfig {
		t.Errorf("NewRepo() failed, expected AppConfig to be %v, got %v", appConfig, repo.App)
	}
}

func TestNewHandlers(t *testing.T) {
	appConfig := &config.AppConfig{}
	repo := handlers.NewRepo(appConfig)
	handlers.NewHandlers(repo)
	if handlers.Repo != repo {
		t.Errorf("NewHandlers() failed, expected Repo to be %v, got %v", repo, handlers.Repo)
	}
}
