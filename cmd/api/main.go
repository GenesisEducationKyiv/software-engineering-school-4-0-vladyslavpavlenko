package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/scheduler"
)

const webPort = 8080

var app config.AppConfig

func main() {
	err := run()
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	dbconn, err := setup(&app)
	if err != nil {
		return err
	}
	defer dbconn.Close()

	s := scheduler.NewCronScheduler()
	schedule := "0 10 * * *" // every day at 10 AM

	_, err = s.ScheduleTask(schedule, func() {
		err = handlers.Repo.NotifySubscribers()
		if err != nil {
			log.Printf("Error notifying subscribers: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to schedule mailer task: %v", err)
	}
	s.Start()
	defer s.Stop()

	log.Printf("Running on port %d", webPort)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", webPort),
		Handler:           routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	// Block until a signal is received
	<-stop

	// Set a deadline
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	log.Println("Shutting down...")
	if err = srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %v", err)
	}

	log.Println("Server has been stopped")
	return nil
}
