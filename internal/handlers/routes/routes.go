package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
	m "github.com/vladyslavpavlenko/genesis-api-project/internal/handlers/middleware"
)

// API sets up the main application routes and middleware for the API.
func API(h *handlers.Handlers) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Heartbeat("/health"))
	mux.Use(middleware.RequestID)
	mux.Use(m.Metrics)

	mux.Route("/api", func(mux chi.Router) {
		mux.Route("/v1", func(mux chi.Router) {
			mux.Get("/rate", h.GetRate)
			mux.Post("/subscribe", h.Subscribe)
			mux.Post("/unsubscribe", h.Unsubscribe)
			mux.Post("/sendEmails", h.SendEmails)
		})
	})

	return mux
}

// Metrics sets up the routes for metrics endpoints.
func Metrics() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Heartbeat("/health"))
	mux.Get("/metrics", handlers.Metrics)

	return mux
}
