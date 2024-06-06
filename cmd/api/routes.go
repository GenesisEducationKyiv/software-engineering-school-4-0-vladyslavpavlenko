package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/"))

	mux.Get("/rate", handlers.Repo.GetRate)
	mux.Post("/subscribe", handlers.Repo.Subscribe)
	mux.Post("/sendEmails", handlers.Repo.SendEmails)

	return mux
}
