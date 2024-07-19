package routes_test

import (
	"testing"

	"github.com/go-chi/chi"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers/routes"
)

func TestRoutes(t *testing.T) {
	mux := routes.Routes()

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Errorf("type is not chi.Mux, but is %T", v)
	}
}
