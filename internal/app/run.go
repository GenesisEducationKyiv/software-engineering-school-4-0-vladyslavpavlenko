package app

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

	"github.com/robfig/cron/v3"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers/routes"
	schedulerpkg "github.com/vladyslavpavlenko/genesis-api-project/internal/scheduler"
)

const (
	webPort  = 8080
	schedule = "0 10 * * *"
)

// scheduler defines an interface for scheduling tasks.
type scheduler interface {
	ScheduleTask(schedule string, task func()) (cron.EntryID, error)
	Start()
	Stop()
}

// dbConnection defines an interface for the database connection.
type dbConnection interface {
	Setup(dsn string) error
	Close() error
	Migrate(models ...any) error
}

// Run initializes the application, sets up the database, schedules email tasks, and starts the
// HTTP server with graceful shutdown.
func Run(appConfig *config.AppConfig) error {
	dbConn, err := setup(appConfig)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dbConn.Close(); closeErr != nil {
			log.Printf("Error closing the database connection: %v\n", closeErr)
		}
	}()

	s := schedulerpkg.NewCronScheduler()
	err = scheduleEmails(s)
	if err != nil {
		return fmt.Errorf("failed to schedule emails: %w", err)
	}
	s.Start()
	defer s.Stop()

	log.Printf("Running on port %d", webPort)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", webPort),
		Handler:           routes.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Graceful shutdown
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

// scheduleEmails uses the provided Scheduler to set up the mailing function.
func scheduleEmails(s scheduler) error {
	_, err := s.ScheduleTask(schedule, func() {
		err := handlers.Repo.NotifySubscribers()
		if err != nil {
			log.Printf("Error notifying subscribers: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to schedule mailing task: %v", err)
	}

	return nil
}
