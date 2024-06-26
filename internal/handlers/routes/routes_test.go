package routes_test

import (
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers/routes"

	"github.com/go-chi/chi"
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
