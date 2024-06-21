package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
)

func Routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	mux.Use(middleware.Heartbeat("/health"))

	mux.Route("/api", func(mux chi.Router) {
		mux.Route("/v1", func(mux chi.Router) {
			mux.Get("/rate", handlers.Repo.GetRate)
			mux.Post("/subscribe", handlers.Repo.Subscribe)
			mux.Post("/sendEmails", handlers.Repo.SendEmails)
		})
	})

	return mux
}
