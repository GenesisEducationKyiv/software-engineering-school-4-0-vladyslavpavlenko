package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	defer app.DB.Close()

	s := scheduler.NewCronScheduler()
	schedule := "0 10 * * *" // every day at 10 AM

	_, err = s.ScheduleTask(schedule, func() {
		err = handlers.Repo.NotifySubscribers()
		if err != nil {
			log.Printf("Error notifying subscribers: %v", err)
		}
	})
	if err != nil {
		log.Printf("Failed to schedule mailer task: %v", err)
		os.Exit(1)
	}
	s.Start()

	log.Printf("Running on port %d", webPort)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", webPort),
		Handler:           routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Printf("HTTP server failed: %v", err)
		os.Exit(1)
	}
}
