package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/scheduler"
)

const webPort = 8080

var app config.AppConfig

func main() {
	err := setup(&app)
	if err != nil {
		log.Fatal(err)
	}

	schedule := "0 10 * * *" // Every day at 10 AM
	_, err = scheduler.ScheduleTask(schedule, handlers.Repo.NotifySubscribers)
	if err != nil {
		log.Fatalf("failed to schedule email task: %v", err)
	}

	log.Printf("Running on port %d", webPort)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", webPort),
		Handler:           routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
