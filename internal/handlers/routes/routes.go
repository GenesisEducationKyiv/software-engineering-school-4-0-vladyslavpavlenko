package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
	m "github.com/vladyslavpavlenko/genesis-api-project/internal/handlers/middleware"
)

func Routes(h *handlers.Handlers) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Heartbeat("/health"))
	mux.Use(m.Metrics)

	mux.Route("/api", func(mux chi.Router) {
		mux.Route("/v1", func(mux chi.Router) {
			mux.Get("/rate", h.GetRate)
			mux.Post("/subscribe", h.Subscribe)
			mux.Post("/unsubscribe", h.Unsubscribe)
			mux.Post("/sendEmails", h.SendEmails)
		})
	})

	mux.Get("/metrics", h.Metrics)

	return mux
}
